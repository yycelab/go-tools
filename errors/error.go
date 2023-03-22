package errors

import (
	"fmt"
)

type ErrType = int

const (
	ERR_PARAM ErrType = iota
	ERR_LOGIC
	ERR_RESOURCE
	ERR_PERMISSION
)

type ParamError struct {
	Field string
	Msg   string
}

func (pe *ParamError) Error() string {
	if len(pe.Field) > 0 {
		return fmt.Sprintf("bad param field(%s):%s", pe.Field, pe.Msg)
	}
	return fmt.Sprintf("bad param %s", pe.Msg)
}

type LogicKind = string

type LogicError struct {
	Msg  string
	Code string
	Kind LogicKind
}

func (le *LogicError) Error() string {
	return fmt.Sprintf("logic error,kind:%s,code:%s,cause:%s", le.Kind, le.Code, le.Msg)
}

type ResourceKind string

const (
	RESOURCE_NETWORK  ResourceKind = "network"
	RESOURCE_DATABASE ResourceKind = "database"
	RESOURCE_CACHE    ResourceKind = "cache"
)

type ErrMessage interface {
	SetMsg(msg string)
	SetMsgIfAbsent(msg string)
}

type ResourceError struct {
	Kind ResourceKind
	Msg  string
}

func (re *ResourceError) Error() string {
	return fmt.Sprintf("resource error,kind:%s,cause:%s", re.Kind, re.Msg)
}

type PermissionError struct {
	Msg     string
	Request string
}

func (pe *PermissionError) Error() string {
	req := pe.Request
	if len(req) == 0 {
		req = "****"
	}
	msg := pe.Msg
	if len(pe.Msg) == 0 {
		msg = "denied,need more permisson"
	}
	return fmt.Sprintf("request:%s,%s", req, msg)
}

func AssertErrorPanic(err error, throw error) {
	if err != nil {
		if throw == nil {
			panic(err)
		}
		panic(throw)
	}
}

type RecoverResult struct {
	HasError   bool
	Err        error
	HttpStatus int
}

func RecoverError(err any) (r *RecoverResult) {
	r = &RecoverResult{HttpStatus: 200}
	if err != nil {
		r.HasError = true
		r.Err, _ = err.(error)
		switch err.(type) {
		case ParamError:
			r.HttpStatus = 405
		case *ParamError:
			r.HttpStatus = 405
		case ResourceError:
			r.HttpStatus = 500
		case *ResourceError:
			r.HttpStatus = 500
		case LogicError:
			r.HttpStatus = 200
		case *LogicError:
			r.HttpStatus = 200
		default:
			r.HttpStatus = 500
		}
	}
	return
}
