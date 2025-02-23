//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/faults
//

package faults

// ErrNotFound creates a basic context for the not found error.
// The error is compatible with `Safe1[string]` and implements `NotFound` interface.
//
//	const errSome = errors.NotFound("key %s is not found")
type ErrNotFound Safe1[string]

// With wraps error into the context.
// The function expands the context with arguments.
//
//	if err := doSomething(); err != nil {
//		return nil, errSome.With(err)
//	}
func (e ErrNotFound) With(err error, key string) error {
	return errNotFound{
		err: Safe1[string](e).With(err, key),
		key: key,
	}
}

func (e ErrNotFound) Error() string { return string(e) }

type errNotFound struct {
	err error
	key string
}

func (e errNotFound) Error() string { return e.err.Error() }

func (e errNotFound) Unwrap() error { return e.err }

func (e errNotFound) NotFound() string { return e.key }
