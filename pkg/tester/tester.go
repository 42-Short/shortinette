package tester

import (
	"fmt"

	"github.com/42-Short/shortinette/internal/config"
)

func Run(configFilePath string) error {
	conf, err := config.GetConfig(configFilePath)
	if err != nil {
		return err
	}
	tests, err := config.GetTests(configFilePath)
	if err != nil {
		return err
	}
	fmt.Println(conf)
	fmt.Println(tests.AssertEq)
	return nil
}
