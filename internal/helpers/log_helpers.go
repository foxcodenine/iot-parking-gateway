package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

// TODO:  to add a buffer to file log, fottow:
// https://chatgpt.com/c/67863b1c-62fc-8013-a1ca-3c0585e901b2

var (
	// Loggers
	infoLog  *log.Logger
	errorLog *log.Logger
	fatalLog *log.Logger

	// Log configuration
	infoLogFile  *os.File
	errorLogFile *os.File
)

func ConfigLogger() {
	// Get log mode from environment variables
	infoLogMode := os.Getenv("INFO_LOG_MODE")   // "console", "file", or "off"
	errorLogMode := os.Getenv("ERROR_LOG_MODE") // "console", "file", or "off"

	// Set up info logger
	switch infoLogMode {
	case "file":
		infoLogFile = setupLogFile("info.log")
		infoLog = log.New(infoLogFile, "INFO\t", log.Ldate|log.Ltime)
	case "off":
		infoLog = log.New(nil, "", 0) // Discard logs
	default:
		infoLog = log.New(os.Stdout, color.New(color.FgBlue).Sprint("INFO\t"), log.Ldate|log.Ltime)
	}

	// Set up error logger
	switch errorLogMode {
	case "file":
		errorLogFile = setupLogFile("error.log")
		errorLog = log.New(errorLogFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
		fatalLog = log.New(errorLogFile, "FATAL\t", log.Ldate|log.Ltime|log.Lshortfile)
	case "off":
		errorLog = log.New(nil, "", 0) // Discard logs
		fatalLog = log.New(nil, "", 0) // Discard logs
	default:
		fmt.Println(os.Getenv("ERROR_LOG_MODE"))
		errorLog = log.New(os.Stderr, color.New(color.FgRed).Sprint("ERROR\t"), log.Ldate|log.Ltime|log.Lshortfile)
		fatalLog = log.New(os.Stderr, color.New(color.FgRed, color.Bold).Sprint("FATAL\t"), log.Ldate|log.Ltime|log.Lshortfile)
	}
}

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
	} else {
		errorLog.Output(callDepth, fmt.Sprintf("%s:\n", msg))
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
	} else {
		fatalLog.Output(callDepth, fmt.Sprintf("%s:\n", msg))
	}

	// Close log files if open
	if infoLogFile != nil {
		infoLogFile.Close()
	}
	if errorLogFile != nil {
		errorLogFile.Close()
	}

	os.Exit(1) // Terminate the application with a non-zero status code
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

// setupLogFile creates or opens a log file for writing.
func setupLogFile(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file %s: %v\n", filename, err)
		os.Exit(1)
	}
	return file
}
