package utils

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

type AppError struct {
	StatusCode int
	Message    string
}

func ResponseSuccess(statusCode int, message string, data interface{}) Response {
	return Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}

func ResponseError(err AppError) Response {
	return Response{
		StatusCode: err.StatusCode,
		Message:    err.Message,
	}
}
