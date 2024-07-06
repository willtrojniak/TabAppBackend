package services

import (
	"encoding/json"
	"net/http"
)

type HTTPError interface {
	Data() interface{}
	Msg() string
	StatusCode() int
	error
}

func writeHttpError(w http.ResponseWriter, e HTTPError) {
	w.WriteHeader(e.StatusCode())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		struct {
			Msg  string      `json:"msg"`
			Data interface{} `json:"data"`
		}{Msg: e.Msg(), Data: e.Data()})

}

type HTTPErrorHandler func(http.ResponseWriter, error)

func HandleHttpError(w http.ResponseWriter, err error) {
	switch err := err.(type) {
	case HTTPError:
		writeHttpError(w, err)
		return
	default:
		writeHttpError(w, NewInternalServiceError(err))
		return
	}
}

type ServiceError struct {
	data interface{}
	msg  string
	code int
	err  error
}

func NewInternalServiceError(err error) *ServiceError {
	return NewServiceError(err, http.StatusInternalServerError, nil)
}

func NewUnauthenticatedServiceError(err error) *ServiceError {
	return NewServiceError(err, http.StatusUnauthorized, nil)
}

func NewUnauthorizedServiceError(err error) *ServiceError {
	return NewServiceError(err, http.StatusForbidden, nil)
}

func NewDataConflictServiceError(err error) *ServiceError {
	return NewServiceError(err, http.StatusConflict, nil)
}

func NewNotFoundServiceError(err error) *ServiceError {
	return NewServiceError(err, http.StatusNotFound, nil)
}

type ValidationError struct {
	Value interface{}
	Error string
}
type ValidationErrors map[string]ValidationError

func NewValidationServiceError(err error, parsingErrors interface{}) *ServiceError {
	return NewServiceError(err, http.StatusBadRequest, parsingErrors)
}

func NewServiceError(err error, code int, data interface{}) *ServiceError {
	return &ServiceError{err: err, code: code, msg: http.StatusText(code), data: data}
}

func (e *ServiceError) StatusCode() int {
	return e.code
}

func (e *ServiceError) Msg() string {
	return e.msg
}

func (e *ServiceError) Data() interface{} {
	return e.data
}

func (e *ServiceError) Error() string {
	return e.err.Error()
}
