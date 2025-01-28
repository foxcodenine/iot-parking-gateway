package validations

import (
	"fmt"
	"net/http"
)

func ValidateEventUpQueryParam(r *http.Request) error {
	event := r.URL.Query().Get("event")
	if event != "up" {
		fmt.Println(fmt.Errorf("invalid or missing event. Only 'event=up' is supported"))
		return fmt.Errorf("invalid or missing event. Only 'event=up' is supported")
	}
	return nil
}
