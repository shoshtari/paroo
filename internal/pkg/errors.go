package pkg

import "errors"

var (
	BadRequestError    = errors.New("bad request")
	UnknownMethodError = errors.New("method is unknown")
)
