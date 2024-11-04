// error_helpers.go
package helpers

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
)

// RespondWithError logs the error, captures file/line, and sends a detailed error response.
func RespondWithError(w http.ResponseWriter, err error, message string, statusCode int) {
	_, file, line, _ := runtime.Caller(1)
	var userError string

	if os.Getenv("DEBUG") == "true" || os.Getenv("DEBUG") == "1" {
		userError = fmt.Sprintf("%s: at %s:%d: \n%v", message, file, line, err)
	} else {
		userError = message
	}

	http.Error(w, userError, statusCode)
}
