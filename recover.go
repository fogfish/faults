package errors

import "errors"

func IsNotFound(err error) bool {
	var e interface{ NotFound() string }

	ok := errors.As(err, &e)
	return ok && e.NotFound() != ""
}
