// log_helpers.go
package helpers

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

var (
	infoLog  = log.New(os.Stdout, color.New(color.FgBlue).Sprint("INFO\t"), log.Ldate|log.Ltime)
	errorLog = log.New(os.Stderr, color.New(color.FgRed).Sprint("ERROR\t"), log.Ldate|log.Ltime|log.Lshortfile)
)

// LogInfo logs an informational message with formatting options.
func LogInfo(format string, args ...interface{}) {
	infoLog.Printf(format, args...)
}

// LogError logs error messages with context including the correct caller location.
func LogError(err error, msg string) {
	// The call depth of 2 usually points to the caller of LogError.
	if err != nil {
		errorLog.Output(2, fmt.Sprintf("%s:\n%v", msg, err))
	}
}

// GetInfoLog returns the shared info log instance.
func GetInfoLog() *log.Logger {
	return infoLog
}

// GetErrorLog returns the shared error log instance.
func GetErrorLog() *log.Logger {
	return errorLog
}
