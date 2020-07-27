package apperror

import (
	"encoding/json"
	"fmt"
)

// GenerateError function
func GenerateError(status int, message string) []byte {
	newError := &AppError{
		Status:  status,
		Message: message,
	}

	errorJSON, err := json.Marshal(newError)

	if err != nil {
		fmt.Print(err)
	}

	return errorJSON
}
