// logger provides logging functionality for various aspects of the application,
// including standard logging for info and error messages, as well as logging for trace files.
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// Loggers for different purposes
var (
	Info    *log.Logger // Info logs general informational messages.
	Warning *log.Logger // Warning logs message with warning prefix messages.
	Error   *log.Logger // Error logs error messages with file and line number information.
	File    *log.Logger // File logs messages to a specified trace file.
)

// GetNewTraceFile generates a new trace file name based on the repository ID and the current timestamp.
//
//   - repoID: the unique identifier for the repository
//
// Returns a string representing the file path for the new trace file.
func GetNewTraceFile(repoID string) string {
	t := time.Now()
	formattedTime := t.Format("20060102_150405")
	return fmt.Sprintf("traces/%s-%s.log", repoID, formattedTime)
}

// InitializeTraceLogger sets up the File logger to write to the specified trace file.
//
//   - filePath: the path to the trace file to which logs should be written
//
// Returns an error if the trace file cannot be created or opened.
func InitializeTraceLogger(filePath string) (err error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	File = log.New(file, "", log.Ldate|log.Ltime)
	originalWriter := File.Writer()
	File.SetOutput(&syncWriter{file: file, writer: originalWriter})
	return nil
}

type syncWriter struct {
	file   *os.File  // The log file
	writer io.Writer // The original writer for the logger
}

// Write writes to the underlying writer and ensures the file is flushed after each write.
func (w *syncWriter) Write(p []byte) (n int, err error) {
	n, err = w.writer.Write(p)
	if err != nil {
		return n, err
	}
	err = w.file.Sync() // Flush to disk after every write
	return n, err
}

// init initializes the Info, Warning, and Error loggers
// with appropriate prefixes and output destinations.
func init() {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	Warning = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
