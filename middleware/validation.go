package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ValidateDataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			err := validateAddEntryData(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func validateAddEntryData(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	var requestData struct {
		Word       string `json:"word"`
		Definition string `json:"definition"`
	}

	err := decoder.Decode(&requestData)
	if err != nil {
		return fmt.Errorf("Invalid JSON payload")
	}

	if len(requestData.Word) < 3 || len(requestData.Word) > 50 {
		return fmt.Errorf("Word length must be between 3 and 50 characters")
	}

	if len(requestData.Definition) < 5 || len(requestData.Definition) > 500 {
		return fmt.Errorf("Definition length must be between 5 and 500 characters")
	}

	return nil
}
