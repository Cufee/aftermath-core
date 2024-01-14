package server

type Response[T any] struct {
	Success bool `json:"success"`
	Data    *T   `json:"data"`

	Error ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Message string `json:"message,omitempty"`
	Context string `json:"context,omitempty"`
}

func NewResponse[T any](data T) Response[T] {
	return Response[T]{Success: true, Data: &data}
}

func NewErrorResponse(message string, context string) Response[any] {
	return Response[any]{Success: false, Error: ResponseError{Message: message, Context: context}}
}

func NewErrorResponseFromError(err error, context ...string) Response[any] {
	if len(context) > 0 {
		return Response[any]{Success: false, Error: ResponseError{Message: err.Error(), Context: context[0]}}
	}
	return Response[any]{Success: false, Error: ResponseError{Message: err.Error()}}
}
