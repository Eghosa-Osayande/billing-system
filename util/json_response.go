package util

type ApiResponseData = any

type SuccessResponseWithData[K ApiResponseData] struct {
	Message string `json:"message"`
	Data    K      `json:"data"`
}
type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

func NewSuccessResponseWithData[K ApiResponseData](msg string, data K) SuccessResponseWithData[K] {
	return SuccessResponseWithData[K]{
		Message: msg,
		Data:    data,
	}
}

func NewSuccessResponse(msg string) SuccessResponse {
	return SuccessResponse{
		Message: msg,
	}
}

func errorMessage(msg string, errors []error) ErrorResponse {
	errStrings := make([]string, len(errors))
	for t := range errors {
		errStrings[t] = errors[t].Error()
	}
	return ErrorResponse{
		Message: msg,
		Errors:  errStrings,
	}
}
