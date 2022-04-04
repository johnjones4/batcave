package api

import "net/http"

type apiError struct {
	status int
	code   int
	parent error
}

const (
	errorCodeInternal      = 1000
	errorCodeStore         = 1001
	errorCodeReqestProcess = 1002
	errorCodeLog           = 1003
	errorCodeExpiredToken  = 2000
)

func wrappedError(e error, code int) apiError {
	return apiError{
		status: http.StatusInternalServerError,
		code:   code,
		parent: e,
	}
}

func (e apiError) Error() string {
	return e.parent.Error()
}

func (e apiError) AppErrCode() int {
	return e.code
}

func (e apiError) HTTPStatus() int {
	return e.status
}
