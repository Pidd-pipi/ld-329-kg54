package errors

import "net/http"

type BusinessError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e BusinessError) Error() string { return e.Message }

func NewBusinessError(status int, code, message string) BusinessError {
	if status == 0 {
		status = http.StatusBadRequest
	}
	return BusinessError{Code: code, Message: message, Status: status}
}
