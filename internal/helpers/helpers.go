package helpers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
)

// ---------------------------------------------------------------------

var (
	infoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
)

// LogInfo logs an informational message with formatting options.
func LogInfo(format string, args ...interface{}) {
	infoLog.Printf(format, args...)
}

// LogError logs error messages with context
func LogError(err error, msg string) {
	errorLog.Printf("%s:\n%v\n", msg, err)
}

// ---------------------------------------------------------------------

// RespondWithError logs the error, captures the file and line number, and sends a detailed error response
func RespondWithError(w http.ResponseWriter, err error, message string, statusCode int) {
	// Capture the file and line number where the error occurred
	_, file, line, _ := runtime.Caller(1) // Caller(1) gets the caller of this function

	// Format the error message with file, line, and provided message
	var userError string

	if os.Getenv("DEBUG") == "true" || os.Getenv("DEBUG") == "1" {
		userError = fmt.Sprintf("%s: at %s:%d: \n%v", message, file, line, err)
	} else {
		userError = fmt.Sprintf("%s", message)
	}

	// Send the error to the client (development only)
	http.Error(w, userError, statusCode)
}
