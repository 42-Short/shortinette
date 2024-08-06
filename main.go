package shortinette

import (
	"os"

	"github.com/42-Short/shortinette/pkg/logger"

	"github.com/42-Short/shortinette/pkg/requirements"
)

func Init() {
	if len(os.Args) == 4 {
		logger.InitializeStandardLoggers(os.Args[2])
	} else {
		logger.InitializeStandardLoggers("")
	}
	if err := requirements.ValidateRequirements(); len(os.Args) != 4 && err != nil {
		logger.Error.Println(err.Error())
		return
	}
}
