// log_helpers.go
package helpers

import (
	"encoding/json"
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
func LogError(err error, msg string, depth ...int) {
	// Set the default depth to 2 if none is provided
	callDepth := 2
	if len(depth) > 0 {
		callDepth = depth[0]
	}
	// The call depth of 2 usually points to the caller of LogError.
	if err != nil {
		errorLog.Output(callDepth, fmt.Sprintf("%s:\n%v", msg, err))
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

// PrettyPrintJSON takes any data and prints it in an indented JSON format
func PrettyPrintJSON(data any) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}
