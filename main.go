package main

import (
	"fmt"
	"os"

	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/R00"
	"github.com/42-Short/shortinette/pkg/requirements"
	Short "github.com/42-Short/shortinette/pkg/short"
	webhook "github.com/42-Short/shortinette/pkg/short/testmodes/webhooktestmode"
)

var DockerFileTemplate = `
FROM debian:latest

RUN apt-get update && apt-get install -y curl build-essential

# Install Go
RUN curl -OL https://go.dev/dl/go1.22.5.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz && \
    rm go1.22.5.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

# Install Rust
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

ENV PATH="/root/.cargo/bin:${PATH}"

WORKDIR /app

CMD ["sh", "-c", "go run . %s %s"]
`

var ModuleOne = map[string]bool{
	"00": true,
	"01": true,
	"02": true,
	"03": true,
	"04": true,
}

var ModulesLookupTable = map[string]interface{}{
	"00": ModuleOne,
}

func dockerExecMode(args []string) error {
	fmt.Println(args)
    module, ok := ModulesLookupTable[args[1]]
    if !ok {
        return fmt.Errorf("module not found")
    }

    moduleMap, ok := module.(map[string]bool)
    if !ok {
        return fmt.Errorf("invalid module type")
    }

    if _, ok := moduleMap[args[2]]; !ok {
        return fmt.Errorf("exercise not found in module")
    }

    fmt.Println("yay")
    return nil
}

func main() {
	logger.InitializeStandardLoggers()
	if len(os.Args) == 3 {
		if err := dockerExecMode(os.Args); err != nil {
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
	short := Short.NewShort("Rust Piscine 1.0", webhook.NewWebhookTestMode())
	config, err := Short.GetConfig()
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	Short.StartModule(*R00.R00(), *config)
	short.TestMode.Run()
	Short.EndModule(*R00.R00(), *config)
}
