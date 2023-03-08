//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/errors
//

package faults

import (
	"errors"
	"time"
)

type Timeout interface{ Timeout() time.Duration }

func IsTimeout(err error, deadline time.Duration) bool {
	var e interface{ Timeout() time.Duration }

	ok := errors.As(err, &e)
	return ok && e.Timeout() >= deadline
}

type NotFound interface{ NotFound() string }

func IsNotFound(err error, key ...string) bool {
	var e interface{ NotFound() string }

	if ok := errors.As(err, &e); !ok {
		return false
	}

	if len(key) == 0 {
		return e.NotFound() != ""
	}

	for _, x := range key {
		if e.NotFound() == x {
			return true
		}
	}

	return false
}

type StatusCode interface{ StatusCode() string }

func IsStatusCode(err error, code ...string) bool {
	var e interface{ StatusCode() string }

	if ok := errors.As(err, &e); !ok {
		return false
	}

	if len(code) == 0 {
		return e.StatusCode() != ""
	}

	for _, x := range code {
		if e.StatusCode() == x {
			return true
		}
	}

	return false
}

type PreConditionFailed interface{ PreConditionFailed() bool }

func IsPreConditionFailed(err error) bool {
	var e interface{ PreConditionFailed() bool }

	ok := errors.As(err, &e)
	return ok && e.PreConditionFailed()
}

type Conflict interface{ Conflict() bool }

func IsConflict(err error) bool {
	var e interface{ Conflict() bool }

	ok := errors.As(err, &e)
	return ok && e.Conflict()
}

type Gone interface{ Gone() bool }

func IsGone(err error) bool {
	var e interface{ Gone() bool }

	ok := errors.As(err, &e)
	return ok && e.Gone()
}

type Issue interface {
	ErrCode() string
	ErrType() string
	ErrInstance() string
	ErrTitle() string
	ErrDetail() string
}
