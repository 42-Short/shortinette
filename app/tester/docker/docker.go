package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/42-Short/shortinette/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type Container struct {
	DockerClient *client.Client
	ID           string
	ExitCode     int64
	Timeout      bool
	Logs         string
}

func NewClient() (*client.Client, error) {
	dockerClient, err := client.NewClientWithOpts(client.WithHost("unix:///var/run/docker.sock"))
	if err != nil {
		return nil, err
	}

	return dockerClient, nil
}

func BuildImage(dockerClient *client.Client, logger *io.Writer) error {
	buildContext, err := archive.TarWithOptions("..", &archive.TarOptions{})
	if err != nil {
		return err
	}

	ctx := context.Background()
	imageBuildResponse, err := dockerClient.ImageBuild(
		ctx,
		buildContext,
		types.ImageBuildOptions{
			Dockerfile: "testenv/Dockerfile",
			Tags:       []string{"shortinette-testenv:latest"},
			Remove:     true,
		},
	)

	if err != nil {
		return err
	}

	defer imageBuildResponse.Body.Close()
	if logger != nil {
		io.Copy(*logger, imageBuildResponse.Body)
	}

	return nil
}

func ContainerCreate(dockerClient *client.Client, command []string) (*Container, error) {
	containerConfig := container.Config{
		Image: "shortinette-testenv",
		Cmd:   command,
	}

	hostConfig := container.HostConfig{}
	networkConfig := network.NetworkingConfig{}
	ctx := context.Background()
	resp, err := dockerClient.ContainerCreate(ctx, &containerConfig, &hostConfig, &networkConfig, nil, "")

	if err != nil {
		return nil, err
	}

	container := Container{
		DockerClient: dockerClient,
		ID:           resp.ID,
	}
	return &container, nil
}

func copyTestExecutableToFolder(exercise config.Exercise, exerciseDirectory string) error {
	src, err := os.Open(exercise.ExecutablePath)

	if err != nil {
		return err
	}
	defer src.Close()

	dstName := filepath.Join(exerciseDirectory, filepath.Base(exercise.ExecutablePath))
	dst, err := os.Create(dstName)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	srcInfo, err := src.Stat()
	if err != nil {
		return err
	}

	if err := os.Chmod(dstName, srcInfo.Mode()); err != nil {
		return err
	}

	return nil
}

func addExecutableToArchive(path string, tarWriter *tar.Writer) error {
	executableInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	executableHeader, err := tar.FileInfoHeader(executableInfo, executableInfo.Name())
	if err != nil {
		return err
	}
	executableHeader.Name = filepath.Base(executableHeader.Name)

	if err := tarWriter.WriteHeader(executableHeader); err != nil {
		return err
	}

	executableData, err := os.Open(path)
	if err != nil {
		return err
	}
	defer executableData.Close()
	if _, err := io.Copy(tarWriter, executableData); err != nil {
		return err
	}

	return nil
}

func createTarArchive(exercise config.Exercise, exerciseDirectory string) (io.Reader, error) {
	var buf bytes.Buffer
	tarWriter := tar.NewWriter(&buf)
	err := filepath.Walk(exerciseDirectory, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, file)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(file[len(exerciseDirectory):])
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(tarWriter, f); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := addExecutableToArchive(exercise.ExecutablePath, tarWriter); err != nil {
		return nil, err
	}

	if err := tarWriter.Close(); err != nil {
		return nil, err
	}
	return &buf, nil
}

func (c *Container) CopyFilesToContainer(exercise config.Exercise, exerciseDirectory string) error {
	tar, err := createTarArchive(exercise, exerciseDirectory)
	if err != nil {
		return err
	}

	ctx := context.Background()
	if err := c.DockerClient.CopyToContainer(ctx, c.ID, "/root", tar, container.CopyToContainerOptions{}); err != nil {
		return err
	}

	return nil
}

func decodeDockerLogs(buffer *bytes.Buffer, logs io.ReadCloser) error {
	tmpbuf := make([]byte, 1024)
	for {
		n, err := logs.Read(tmpbuf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("unable to read docker logs: %s", err)
		}
		if n == 0 {
			break
		}

		logOutput := tmpbuf[:n]
		for len(logOutput) > 0 {
			if len(logOutput) < 8 {
				break
			}

			header := logOutput[:4]
			length := int(header[1])<<16 | int(header[2])<<8 | int(header[3])

			if length > len(logOutput)-4 {
				break
			}

			msg := logOutput[4 : 4+length]
			if _, err := buffer.Write(msg); err != nil {
				return fmt.Errorf("unable to write logs into buffer: %s", err)
			}

			logOutput = logOutput[4+length:]
		}
	}
	return nil
}

func (c *Container) Kill() error {
	ctx := context.Background()

	if err := c.DockerClient.ContainerKill(ctx, c.ID, "SIGKILL"); err != nil {
		return err
	}
	return c.DockerClient.ContainerRemove(ctx, c.ID, container.RemoveOptions{Force: true})
}

func (c *Container) wait(timeout time.Duration) error {
	ctx := context.Background()

	time.AfterFunc(timeout, func() {
		if err := c.DockerClient.ContainerStop(ctx, c.ID, container.StopOptions{}); err == nil {
			c.Timeout = true
		}
	})

	statusCh, errCh := c.DockerClient.ContainerWait(ctx, c.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		c.Kill()
		return fmt.Errorf("error waiting for container: %s", err)
	case status := <-statusCh:
		c.ExitCode = status.StatusCode
	}

	var buf bytes.Buffer
	logs, err := c.DockerClient.ContainerLogs(ctx, c.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return fmt.Errorf("error getting logs from container: %s", err)
	}

	defer logs.Close()
	if err := decodeDockerLogs(&buf, logs); err != nil {
		return err
	}

	c.Logs = buf.String()
	return nil
}

func (c *Container) Exec(timeout time.Duration) error {
	ctx := context.Background()

	if err := c.DockerClient.ContainerStart(ctx, c.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("error starting Docker container: %s", err)
	}

	err := c.wait(timeout)
	c.DockerClient.ContainerRemove(ctx, c.ID, container.RemoveOptions{Force: true})
	return err
}
