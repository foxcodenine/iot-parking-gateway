// log_helpers.go
package helpers

import (
	"log"
	"os"
)

var (
	infoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
)

// LogInfo logs an informational message with formatting options.
func LogInfo(format string, args ...interface{}) {
	infoLog.Printf(format, args...)
}

// LogError logs error messages with context.
func LogError(err error, msg string) {
	errorLog.Printf("%s:\n%v\n", msg, err)
}

// GetInfoLog returns the shared info log instance.
func GetInfoLog() *log.Logger {
	return infoLog
}

// GetErrorLog returns the shared error log instance.
func GetErrorLog() *log.Logger {
	return errorLog
}
