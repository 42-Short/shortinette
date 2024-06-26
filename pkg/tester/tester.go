package tester

import (
	"fmt"

	"github.com/42-Short/shortinette/internal/config"
	"github.com/42-Short/shortinette/pkg/functioncheck"
)

func Run(configFilePath string) error {
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

	fmt.Println(conf.Ex00.Tests)
	return nil
}
