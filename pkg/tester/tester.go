package tester

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/42-Short/shortinette/internal/config"
	"github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/pkg/functioncheck"
	"github.com/42-Short/shortinette/pkg/git"
)

func compileStudentCode(codeDirectory string, turnInDirectory string, turnInFile string) error {
	parentDirectory := fmt.Sprintf("./%s/%s/", codeDirectory, turnInDirectory)

	if _, err := os.Stat(fmt.Sprintf("%s/Cargo.toml", parentDirectory)); os.IsNotExist(err) {
		cmd := exec.Command("rustc", turnInFile)
		cmd.Dir = parentDirectory

		output, err := cmd.CombinedOutput()
		if err != nil {
			return errors.NewSubmissionError(errors.ErrInvalidCompilation, string(output))
		}
	}
	return nil
}

func Run(configFilePath string, studentLogin string, codeDirectory string) error {
	defer func() {
		os.RemoveAll(codeDirectory)
	}()
	allowedItems, err := config.GetAllowedItems(configFilePath)
	if err != nil {
		return err
	}
	conf, err := config.GetConfig(configFilePath)
	if err != nil {
		return err
	}
	if err = functioncheck.Execute(allowedItems, conf.Ex00.TurnInDirectory, conf.Ex00.TurnInFile); err != nil {
		return err
	}
	if err = git.Get(fmt.Sprintf("https://github.com/42-Short/%s.git", studentLogin), codeDirectory); err != nil {
		return err
	}
	if err = compileStudentCode(codeDirectory, conf.Ex00.TurnInDirectory, conf.Ex00.TurnInFile); err != nil {
		return err
	}
	return nil
}
