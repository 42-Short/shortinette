package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	Info     *log.Logger
	Error    *log.Logger
	File     *log.Logger
	Exercise *log.Logger
)

func GetNewTraceFile(repoID string) string {
	t := time.Now()
	formattedTime := t.Format("20060102_150405")
	return fmt.Sprintf("traces/%s-%s.log", repoID, formattedTime)

}

func InitializeTraceLogger(filePath string) (err error) {
	if err = os.Mkdir("traces", 0655); err != nil && !os.IsExist(err) {
		return err
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	File = log.New(file, "", log.Ldate|log.Ltime)
	return nil
}

func InitializeStandardLoggers(exerciseId string) {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	Exercise = log.New(os.Stderr, fmt.Sprintf("EXERCISE %s: ", exerciseId), log.Ldate|log.Ltime|log.Lshortfile)
}
