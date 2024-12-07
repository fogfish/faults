//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/errors
//

// Package errors provides type safe constructs to annotate Golang errors
// with the context and handle opaque errors without boilerplate
// https://tech.fog.fish/2022/07/05/assert-golang-errors-for-behavior.html#opaque-errors
//
// It solves a problem of annotate errors with the context so that
// consequent debugging of opaque error handling becomes an easy job.
// Instead of using `fmt.Errorf` to include the execution context to the error,
// it defines a type safe wrapping of errors.
package faults

import (
	"fmt"
	"runtime"
)

// Type creates a basic context for the error. The context produces an error like
// `[function line] text defined by context: original error`
//
//	const errSome = errors.Type("something is failed")
type Type string

// With wraps error into the context.
// The function expands the context with arguments.
//
//	if err := doSomething(); err != nil {
//		return nil, errSome.With(err)
//	}
func (e Type) With(err error, args ...any) error {
	var (
		name string
		line int
	)

	if pc, _, ln, ok := runtime.Caller(1); ok {
		name = runtime.FuncForPC(pc).Name()
		line = ln
	}

	msg := string(e)
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	return fmt.Errorf("[%s %d] "+msg+": %w", name, line, err)
}

// Deprecated: Use With
func (e Type) New(err error, args ...any) error {
	return e.With(err, args...)
}

func (e Type) Error() string { return string(e) }

// Fast creates a basic context for the error but skips usage of runtime package.
//
//	const errSome = errors.Fast("something is failed")
type Fast string

// With wraps error into the context.
// The function expands the context with arguments.
//
//	if err := doSomething(); err != nil {
//		return nil, errSome.With(err)
//	}
func (e Fast) With(err error, args ...any) error {
	msg := string(e)
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	return fmt.Errorf(msg+": %w", err)
}

// Deprecated: Use With
func (e Fast) New(err error, args ...any) error {
	return e.With(err, args...)
}

func (e Fast) Error() string { return string(e) }

// Safe1 creates an error context with 1 argument
//
//	const errSome = errors.Safe1[string]("something is failed %s")
type Safe1[A any] string

// With wraps error into the context.
// The function expands the context with arguments.
//
//	if err := doSomething(); err != nil {
//		return nil, errSome.With(err, "foo")
//	}
func (safe Safe1[A]) With(err error, a A) error {
	var (
		name string
		line int
	)

	if pc, _, ln, ok := runtime.Caller(1); ok {
		name = runtime.FuncForPC(pc).Name()
		line = ln
	}

	return fmt.Errorf("[%s %d] "+string(safe)+": %w", name, line, a, err)
}

// Deprecated: Use With
func (safe Safe1[A]) New(err error, a A) error {
	return safe.With(err, a)
}

// Safe2 creates an error context with 2 argument
type Safe2[A, B any] string

// With wraps error into the context.
func (safe Safe2[A, B]) With(err error, a A, b B) error {
	var (
		name string
		line int
	)

	if pc, _, ln, ok := runtime.Caller(1); ok {
		name = runtime.FuncForPC(pc).Name()
		line = ln
	}

	return fmt.Errorf("[%s %d] "+string(safe)+": %w", name, line, a, b, err)
}

// Deprecated: Use With
func (safe Safe2[A, B]) New(err error, a A, b B) error {
	return safe.With(err, a, b)
}

// Safe3 creates an error context with 3 argument
type Safe3[A, B, C any] string

// With wraps error into the context.
func (safe Safe3[A, B, C]) With(err error, a A, b B, c C) error {
	var (
		name string
		line int
	)

	if pc, _, ln, ok := runtime.Caller(1); ok {
		name = runtime.FuncForPC(pc).Name()
		line = ln
	}

	return fmt.Errorf("[%s %d] "+string(safe)+": %w", name, line, a, b, c, err)
}

// Deprecated: Use With
func (safe Safe3[A, B, C]) New(err error, a A, b B, c C) error {
	return safe.With(err, a, b, c)
}

// Safe4 creates an error context with 4 argument
type Safe4[A, B, C, D any] string

// With wraps error into the context.
func (safe Safe4[A, B, C, D]) With(err error, a A, b B, c C, d D) error {
	var (
		name string
		line int
	)

	if pc, _, ln, ok := runtime.Caller(1); ok {
		name = runtime.FuncForPC(pc).Name()
		line = ln
	}

	return fmt.Errorf("[%s %d] "+string(safe)+": %w", name, line, a, b, c, d, err)
}

// Deprecated: Use With
func (safe Safe4[A, B, C, D]) New(err error, a A, b B, c C, d D) error {
	return safe.With(err, a, b, c, d)
}

// Safe5 creates an error context with 5 argument
type Safe5[A, B, C, D, E any] string

// With wraps error into the context.
func (safe Safe5[A, B, C, D, E]) With(err error, a A, b B, c C, d D, e E) error {
	var (
		name string
		line int
	)

	if pc, _, ln, ok := runtime.Caller(1); ok {
		name = runtime.FuncForPC(pc).Name()
		line = ln
	}

	return fmt.Errorf("[%s %d] "+string(safe)+": %w", name, line, a, b, c, d, e, err)
}

// Deprecated: Use With
func (safe Safe5[A, B, C, D, E]) New(err error, a A, b B, c C, d D, e E) error {
	return safe.With(err, a, b, c, d, e)
}
