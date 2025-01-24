// error_helpers.go
package helpers

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/fatih/color"
)

// RespondWithError logs the error, captures file/line, and sends a detailed error response.
func RespondWithError(w http.ResponseWriter, err error, message string, statusCode int, depth ...int) {
	_, file, line, _ := runtime.Caller(1)
	var userError string

	// Set the default depth to 2 if none is provided
	callDepth := 3
	if len(depth) > 0 {
		callDepth = depth[0]
	}

	if os.Getenv("DEBUG") == "true" || os.Getenv("DEBUG") == "1" {
		userError = fmt.Sprintf("%s: at %s:%d: \n%v", message, file, line, err)
		LogError(err, message, callDepth)
	} else {
		userError = message
		LogError(err, message, callDepth)
	}

	http.Error(w, userError, statusCode)
}

// Helper function to wrap error with colored file and line context
func WrapError(err error) error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("%w", err)
	}

	fileLineInfo := fmt.Sprintf("[%s:%d]", file, line)

	if os.Getenv("ERROR_LOG_MODE") != "file" {
		// Use color to make file and line number red
		red := color.New(color.FgRed).SprintFunc()
		fileLineInfo = red(fmt.Sprintf("[%s:%d]", file, line))
	}

	return fmt.Errorf("%s: %w", fileLineInfo, err)
}
