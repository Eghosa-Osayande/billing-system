package util

import "fmt"

type ApiError struct {
	Field   *string `json:"field"`
	Message string  `json:"message"`
}

func (apiError ApiError) Error() string {
	return apiError.Message
}

type ApiErrorList []error

func (apiError ApiErrorList) Error() string {
	errorString := ""
	for index, err := range apiError {
		if index > 0 {
			errorString += ", "
		}
		errorString += fmt.Sprintf("%v", err.Error())
	}
	return errorString
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  []error     `json:"errors"`
}

func SuccessMessage(msg string, data interface{}) ApiResponse {
	return ApiResponse{
		Success: true,
		Message: msg,
		Data:    data,
		Errors:  nil,
	}
}

func errorMessage(msg string, errors []error) ApiResponse {
	return ApiResponse{
		Success: false,
		Message: msg,
		Data:    nil,
		Errors:  errors,
	}
}
