package pkg

import "net/http"

type ApiResponse[T any] struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	Data            T      `json:"data"`
}

func BuildResponse[T any](responseStatusCode int, data T) ApiResponse[T] {
	return BuildResponse_(responseStatusCode, http.StatusText(responseStatusCode), data)
}

func BuildResponse_[T any](code int, message string, data T) ApiResponse[T] {
	return ApiResponse[T]{
		ResponseCode:    code,
		ResponseMessage: message,
		Data:            data,
	}
}

// for Error Response
type ApiWithoutResponse struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
}

func BuildWithoutResponse(responseStatusCode int, msg string) ApiWithoutResponse {
	return BuildWithoutResponse_(responseStatusCode, msg)
}

func BuildWithoutResponse_(code int, message string) ApiWithoutResponse {
	return ApiWithoutResponse{
		ResponseCode:    code,
		ResponseMessage: message,
	}
}
