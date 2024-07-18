package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/R00"
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
	"github.com/42-Short/shortinette/pkg/requirements"
	Short "github.com/42-Short/shortinette/pkg/short"
	webhook "github.com/42-Short/shortinette/pkg/short/testmodes/webhooktestmode"
)

func dockerExecMode(args []string, short Short.Short) error {
	exercise, ok := short.Modules[args[1]].Exercises[args[2]]
	if !ok {
		return fmt.Errorf("could not find exercise")
	}
	if err := logger.InitializeTraceLogger(args[3]); err != nil {
		return err
	}
	result := exercise.Run()
	logger.File.Printf("[MOD%s][EX%s]: %s", args[1], args[2], result.Output)
	if result.Passed {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
	return nil
}

func printCircular(lines []string) {
	for range len(lines) {
		fmt.Printf("\033[A")
		fmt.Printf("\033[K")
	}
	for _, line := range lines {
		fmt.Printf("\033[90m%s\033[0m\n", line)
	}
}

func captureOutput(r io.Reader, done chan<- bool) []string {
	scanner := bufio.NewScanner(r)
	buffer := make([]string, 5)
	var lines []string
	idx := 0

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

		buffer[idx%5] = line
		idx++
		printCircular(buffer)
	}

	done <- true
	return lines
}

func runCommandWithLimitedOutput(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	stdoutDone := make(chan bool)
	stderrDone := make(chan bool)

	go captureOutput(stdoutPipe, stdoutDone)
	go captureOutput(stderrPipe, stderrDone)

	<-stdoutDone
	<-stderrDone

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func buildDockerTestEnvironment() error {
	logger.Info.Println("pre-built testenv image not found, building from ./Dockerfile...")
	err := runCommandWithLimitedOutput("sh", "-c", "docker image ls | grep testenv")
	fmt.Printf("\033[S")
	fmt.Printf("\033[S")
	fmt.Printf("\033[S")
	if err != nil {
		err = runCommandWithLimitedOutput("docker", "build", "-t", "testenv", ".")
		if err != nil {
			logger.Info.Println("in order to compile and test submissions in a safe environment, you will need to add a Dockerfile with all necessary dependencies to the root of your project - see http://github.com/42-Short/shortinette/.github/docs for more details")
			return err
		}
	}
	return nil
}

func main() {
	logger.InitializeStandardLoggers()

	if err := buildDockerTestEnvironment(); err != nil {
		logger.Error.Println(err)
		return
	}

	short := Short.NewShort("Rust Piscine 1.0", map[string]Module.Module{"00": *R00.R00()}, webhook.NewWebhookTestMode())
	if len(os.Args) == 4 {
		if err := dockerExecMode(os.Args, short); err != nil {
			logger.Error.Println(err)
			return
		}
		return
	} else if len(os.Args) != 1 {
		logger.Error.Println("invalid number of arguments")
		return
	}
	if err := requirements.ValidateRequirements(); err != nil {
		logger.Error.Println(err.Error())
		return
	}
	config, err := Short.GetConfig()
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	Short.StartModule(*R00.R00(), *config)
	short.TestMode.Run()
	Short.EndModule(*R00.R00(), *config)
}
