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
	fatalLog = log.New(os.Stderr, color.New(color.FgRed, color.Bold).Sprint("FATAL\t"), log.Ldate|log.Ltime|log.Lshortfile)
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
	if err != nil {
		errorLog.Output(callDepth, fmt.Sprintf("%s:\n%v", msg, err))
	}
}

// LogFatal logs critical errors and terminates the application.
func LogFatal(err error, msg string, depth ...int) {
	callDepth := 2 // Default call depth
	if len(depth) > 0 {
		callDepth = depth[0]
	}
	if err != nil {
		fatalLog.Output(callDepth, fmt.Sprintf("%s:\n%v", msg, err))
	}
	os.Exit(1) // Terminate the application with a non-zero status code
}

// GetInfoLog returns the shared info log instance.
func GetInfoLog() *log.Logger {
	return infoLog
}

// GetErrorLog returns the shared error log instance.
func GetErrorLog() *log.Logger {
	return errorLog
}

// GetErrorLog returns the shared error log instance.
func GetFatalLog() *log.Logger {
	return fatalLog
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
