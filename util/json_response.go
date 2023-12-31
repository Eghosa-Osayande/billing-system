package util

type ApiResponseData = any

type SuccessResponse[K ApiResponseData] struct {
	Message string `json:"message"`
	Data    K      `json:"data"`
}

type ErrorResponse struct {
	Message string  `json:"message"`
	Errors  []string `json:"errors"`
}

func SuccessMessage[K ApiResponseData](msg string, data K) SuccessResponse[K] {
	return SuccessResponse[K]{

		Message: msg,
		Data:    data,
	}
}

func errorMessage(msg string, errors []error) ErrorResponse {
	errStrings := make([]string, len(errors))
	for t:= range errors {
		errStrings[t] = errors[t].Error()
	}
	return ErrorResponse{
		Message: msg,
		Errors: errStrings,
	}
}
