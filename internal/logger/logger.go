package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	Info  *log.Logger
	Error *log.Logger
	File  *log.Logger
)

func GetNewTraceFile(repoID string) string {
	t := time.Now()
	formattedTime := t.Format("20060102_150405")
	return fmt.Sprintf("traces/%s-%s.log", repoID, formattedTime)

}

func InitializeTraceLogger(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	File = log.New(file, "", log.Ldate|log.Ltime)
	return nil
}

func InitializeStandardLoggers() {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
