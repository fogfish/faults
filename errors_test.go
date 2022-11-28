//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/errors
//

package errors_test

import (
	"fmt"
	"testing"

	"github.com/fogfish/errors"
)

func TestType(t *testing.T) {
	const errA = errors.Type("a")

	if errA.New(err).Error() != "[github.com/fogfish/errors_test.TestType 21] a: just error" {
		t.Errorf("failed: %s", errA.New(err))
	}

	const errB = errors.Type("b %s")

	if errB.New(err, "b").Error() != "[github.com/fogfish/errors_test.TestType 27] b b: just error" {
		t.Errorf("failed: %s", errB.New(err, "b"))
	}
}

func TestFast(t *testing.T) {
	const errA = errors.Fast("a")

	if errA.New(err).Error() != "a: just error" {
		t.Errorf("failed: %s", errA.New(err))
	}

	const errB = errors.Fast("b %s")

	if errB.New(err, "b").Error() != "b b: just error" {
		t.Errorf("failed: %s", errB.New(err, "b"))
	}
}

func TestSafe(t *testing.T) {
	const errA = errors.Safe1[string]("a %s")

	if errA.New(err, "a").Error() != "[github.com/fogfish/errors_test.TestSafe 49] a a: just error" {
		t.Errorf("failed: %s", errA.New(err, "a"))
	}

	const errB = errors.Safe2[string, string]("a %s %s")

	if errB.New(err, "a", "b").Error() != "[github.com/fogfish/errors_test.TestSafe 55] a a b: just error" {
		t.Errorf("failed: %s", errB.New(err, "a", "b"))
	}

	const errC = errors.Safe3[string, string, string]("a %s %s %s")

	if errC.New(err, "a", "b", "c").Error() != "[github.com/fogfish/errors_test.TestSafe 61] a a b c: just error" {
		t.Errorf("failed: %s", errC.New(err, "a", "b", "c"))
	}

	const errD = errors.Safe4[string, string, string, string]("a %s %s %s %s")

	if errD.New(err, "a", "b", "c", "d").Error() != "[github.com/fogfish/errors_test.TestSafe 67] a a b c d: just error" {
		t.Errorf("failed: %s", errD.New(err, "a", "b", "c", "d"))
	}

	const errE = errors.Safe5[string, string, string, string, string]("a %s %s %s %s %s")

	if errE.New(err, "a", "b", "c", "d", "e").Error() != "[github.com/fogfish/errors_test.TestSafe 73] a a b c d e: just error" {
		t.Errorf("failed: %s", errE.New(err, "a", "b", "c", "d", "e"))
	}
}

// ------------------------------------------------------------------------------
//
// # Benchmark
//
// ------------------------------------------------------------------------------
var (
	err = fmt.Errorf("just error")
	glo error
)

const (
	errFast = errors.Fast("error fast")
	errType = errors.Type("error type")
	errSafe = errors.Safe1[string]("error %s")
)

func failStdr() error { return fmt.Errorf("error type: %w", err) }
func failFast() error { return errFast.New(err) }
func failType() error { return errType.New(err) }
func failSafe() error { return errSafe.New(err, "safe") }

func BenchmarkStd(b *testing.B) {
	var err error

	for n := 0; n < b.N; n++ {
		err = failStdr()
	}

	glo = err
}

func BenchmarkFast(b *testing.B) {
	var err error

	for n := 0; n < b.N; n++ {
		err = failFast()
	}

	glo = err
}

func BenchmarkType(b *testing.B) {
	var err error

	for n := 0; n < b.N; n++ {
		err = failType()
	}

	glo = err
}

func BenchmarkSafe(b *testing.B) {
	var err error

	for n := 0; n < b.N; n++ {
		err = failSafe()
	}

	glo = err
}
