package pkg

import "errors"

var (
	BadRequestError     = errors.New("bad request")
	UnknownMethodError  = errors.New("method is unknown")
	InternalError       = errors.New("internal error")
	NotImplementedError = errors.New("unimplemented")
	NotFoundError       = errors.New("not found")
)
