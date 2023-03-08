package ecode

import (
	"fmt"

	"github.com/pkg/errors"
	gcode "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var codes = map[int]struct{}{}

func New(code int, msg string) Error {
	if code < 1000 {
		panic("error code must be greater than 1000")
	}
	return add(code, msg)
}

func add(code int, msg string) Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", code))
	}
	codes[code] = struct{}{}
	return Error{
		code: code,
		msg:  msg,
	}
}

type Errors interface {
	Error() string
	Code() int
	Message() string
	Details() []interface{}
	Equal(error) bool
	Reload(string) Error
	GRPCStatus() *status.Status
}

type Error struct {
	code       int
	msg        string
	grpcStatus *status.Status
}

func (e Error) Error() string {
	return e.msg
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Message() string {
	return e.msg
}

func (e Error) Details() []interface{} {
	return nil
}

func (e Error) Equal(err error) bool {
	return Equal(err, e)
}

func (e Error) Reload(s string) Error {
	e.msg = s
	return e
}

func (e Error) GRPCStatus() *status.Status {
	return status.New(gcode.Code(e.Code()), e.msg)
}

func Equal(err error, e Error) bool {
	return Cause(err).Code() == e.Code()
}

func Cause(err error) Errors {
	if err == nil {
		return Success
	}
	if ec, ok := errors.Cause(err).(Errors); ok {
		return ec
	}
	if st, ok := status.FromError(err); ok {
		e := Error{code: int(st.Code()), msg: st.Message()}
		if ec, ok := errors.Cause(e).(Errors); ok {
			return ec
		}
	}
	if err.Error() == "" {
		return Success
	}
	return Error{code: 500, msg: err.Error()}
}
