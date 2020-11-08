package jsonrpc2

import "errors"

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// some built-in errors
var (
	ErrParseFailed    = errors.New("parse error")
	ErrInternalError  = errors.New("invalid params")
	ErrInvalidParams  = errors.New("invalid params")
	ErrMethodNotFound = errors.New("method not found")
	ErrInvalidRequest = errors.New("invalid request")
)

// if is built-in errors, code can be 0 or any.
func NewError(code int, err error, moreDataParams ...interface{}) *Error {
	switch err {
	case ErrParseFailed:
		return newError(-32700, err.Error(), moreDataParams...)
	case ErrInternalError:
		return newError(-32603, err.Error(), moreDataParams...)
	case ErrInvalidParams:
		return newError(-32602, err.Error(), moreDataParams...)
	case ErrMethodNotFound:
		return newError(-32601, err.Error(), moreDataParams...)
	case ErrInvalidRequest:
		return newError(-32600, err.Error(), moreDataParams...)
	default:
		return newError(code, err.Error(), moreDataParams...)
	}
}

func newError(code int, msg string, data ...interface{}) *Error {
	var errData interface{}
	switch len(data) {
	case 0:
		errData = nil
	case 1:
		errData = data[0]
	default:
		errData = data
	}

	return &Error{
		Code:    code,
		Message: msg,
		Data:    errData,
	}
}
