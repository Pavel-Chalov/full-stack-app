package lib

type WebError struct {
	Status int
	Error  string
}

func NewWebError(status int, text string) *WebError {
	return &WebError{Status: status, Error: text}
}

func Conflict(text string) *WebError {
	return &WebError{Status: 409, Error: text}
}

func NotFound(text string) *WebError {
	return &WebError{Status: 404, Error: text}
}

func Unprocessable(text string) *WebError {
	return &WebError{Status: 422, Error: text}
}

func Forbidden(text string) *WebError {
	return &WebError{Status: 403, Error: text}
}

func Unauthorized(text string) *WebError {
	return &WebError{Status: 401, Error: text}
}

func BadRequest(text string) *WebError {
	return &WebError{Status: 400, Error: text}
}

func ServerError(text string) *WebError {
	return &WebError{Status: 500, Error: text}
}
