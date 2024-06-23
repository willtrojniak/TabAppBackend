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
  w.WriteHeader(e.StatusCode());
  json.NewEncoder(w).Encode(
    struct {
      Msg string `json:"msg"`
      Data interface{} `json:"data"`
    }{Msg: e.Msg(), Data: e.Data()});
  
}

func HandleHttpError(w http.ResponseWriter, err error) {
  switch err := err.(type) {
  case HTTPError:
    writeHttpError(w, err);
    return;
  default:
    writeHttpError(w, NewInternalServiceError(err));
    return;
  }
}

type ServiceError struct {
  data interface{}
  msg string 
  code int
  err error
}

func NewInternalServiceError(err error) *ServiceError {
  return NewServiceError(err, http.StatusInternalServerError, "Internal server error.", nil);
}


func NewUnauthorizedServiceError(err error) *ServiceError {
  return NewServiceError(err, http.StatusUnauthorized, "Unauthorized.", nil);
}

func NewServiceError(err error, code int, msg string, data interface{}) *ServiceError {
  return &ServiceError{err: err, code: code, msg: msg, data: data};
}

func (e *ServiceError) StatusCode() int {
  return e.code;
}

func (e *ServiceError) Msg() string {
  return e.msg;
}

func (e *ServiceError) Data() interface{} {
  return e.data;
}

func (e *ServiceError) Error() string {
  return e.Error();
}
