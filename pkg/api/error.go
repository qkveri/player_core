package api

import "fmt"

type NoInternetError struct {
	Err error
}

func (e *NoInternetError) Error() string {
	return fmt.Sprintf("no internet connection: %v", e.Err)
}

type RequestError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("API Request Error: {statusCode: %d, message: %s}",
		e.StatusCode,
		e.Message)
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("API Unauthorized: %s", e.Message)
}

type ValidationError struct {
	StatusCode      int                 `json:"statusCode"`
	Message         string              `json:"message"`
	ValidationFails map[string][]string `json:"validationFails"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("API Validation Error: {statusCode: %d, message: %s, fails: %v}",
		e.StatusCode,
		e.Message,
		e.ValidationFails)
}
