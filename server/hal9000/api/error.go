package api

import (
	"log"
	"net/http"
)

type apiError struct {
	status int
	code   int
	parent error
}

func wrappedError(e error, code int) apiError {
	log.Println(e)
	return apiError{
		status: http.StatusInternalServerError,
		code:   code,
		parent: e,
	}
}

func (e apiError) Error() string {
	return "Sorry. Something went wrong."
}

func (e apiError) AppErrCode() int {
	return e.code
}

func (e apiError) HTTPStatus() int {
	return e.status
}
