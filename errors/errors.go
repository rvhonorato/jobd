// Package errors provides the error handling for the jobd application
package errors

import "net/http"

type RestErr struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   "internal_server_error",
	}
}

func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   "bad_request",
	}

}

func NewConflictError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusConflict,
		Error:   "conflict",
	}
}

func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   "not_found",
	}
}

func NewStatusUnauthorized(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusUnauthorized,
		Error:   "unauthorized",
	}
}

func NewStatusAccepted(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusAccepted,
		Error:   "no_content",
	}
}
