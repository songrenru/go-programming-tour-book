package errcode

import (
	"fmt"
	"net/http"
)

type Error struct {
	code int `json:"code"`
	msg string `json:"msg"`
	details []string `json:"details"`
}

var codes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已经存在，请更换一个", code))
	}
	codes[code] = msg
	return &Error{
		code: code,
		msg: msg,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("错误码：%d, 错误信息：%s", e.code, e.msg)
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Msg() string {
	return e.msg
}

func (e Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg, args)
}

func (e Error) Detail() []string {
	return e.details
}

func (e *Error) WriteDetails(details ...string) *Error {
	newErr := *e
	newErr.details = []string{}
	for _, detail := range details {
		newErr.details = append(newErr.details, detail)
	}

	return &newErr
}

func (e Error) StatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case InvalidParams.Code():
		return http.StatusBadRequest
	case NotFound.Code():
		return http.StatusNotFound
	case UnauthorizedTokenTimeout.Code():
		return http.StatusUnauthorized
	case TooManyRequests.Code():
		return http.StatusTooManyRequests
	case ServerError.Code():
	case UnauthorizedAuthNotExist.Code():
	case UnauthorizedTokenError.Code():
	case UnauthorizedTokenGenerate.Code():
	}
	return http.StatusInternalServerError
}

