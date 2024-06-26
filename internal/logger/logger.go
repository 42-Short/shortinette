package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/42-Short/shortinette/internal/errors"
)

var (
	Info  *log.Logger
	Error *log.Logger
	File  *log.Logger
)

func InitializeTraceLogger(repoId string) error {
	t := time.Now()
	formattedTime := t.Format("20060102_150405")
	fileName := fmt.Sprintf("logs/%s-%s.log", repoId, formattedTime)

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return errors.NewInternalError(errors.ErrInternal, err.Error())
	}
	File = log.New(file, "", log.Ldate|log.Ltime)
	return nil
}

func InitializeStandardLoggers() {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
