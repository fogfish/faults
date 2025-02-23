//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/faults
//

package faults_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/fogfish/faults"
)

func check(t *testing.T, err error, re string) {
	t.Helper()

	r, exx := regexp.Compile(re)
	if exx != nil {
		t.Errorf("invalid regexp: %s", re)
		return
	}

	if !r.MatchString(err.Error()) {
		t.Errorf("failed: %v, expect: %s", err, re)
	}
}

func checkIs(t *testing.T, input, target error) {
	t.Helper()

	if !errors.Is(input, target) {
		t.Errorf("errors.Is is failed, input: %s, target: %s", input, target)
	}
}

func TestType(t *testing.T) {
	const errA = faults.Type("a")
	check(t, errA, "^a$")
	check(t, errA.With(err), "^\\[github.com/fogfish/faults_test.TestType [0-9]+\\] a: just error$")

	const errB = faults.Type("b %s")
	check(t, errB.With(err, "b"), "^\\[github.com/fogfish/faults_test.TestType [0-9]+\\] b b: just error$")

	checkIs(t, errA, errA)
	checkIs(t, errA.With(err), errA)
	checkIs(t, errA.With(err), err)
	checkIs(t, errA.With(errB), errB)
	checkIs(t, errA.With(errB.With(err)), err)
}

func TestFast(t *testing.T) {
	const errA = faults.Fast("a")
	check(t, errA, "a")
	check(t, errA.With(err), "a: just error")

	const errB = faults.Fast("b %s")
	check(t, errB, "b %s")
	check(t, errB.With(err, "b"), "b b: just error")

	checkIs(t, errA, errA)
	checkIs(t, errA.With(err), errA)
	checkIs(t, errA.With(err), err)
	checkIs(t, errA.With(errB), errB)
	checkIs(t, errA.With(errB.With(err)), err)
}

func TestSafe(t *testing.T) {
	const errA = faults.Safe1[string]("a %s")
	check(t, errA, "a %s")
	check(t, errA.With(err, "a"), "^\\[github.com/fogfish/faults_test.TestSafe [0-9]+\\] a a: just error$")

	checkIs(t, errA, errA)
	checkIs(t, errA.With(err, "a"), errA)
	checkIs(t, errA.With(err, "a"), err)

	const errB = faults.Safe2[string, string]("a %s %s")
	check(t, errB, "a %s %s")
	check(t, errB.With(err, "a", "b"), "^\\[github.com/fogfish/faults_test.TestSafe [0-9]+\\] a a b: just error$")

	checkIs(t, errB, errB)
	checkIs(t, errB.With(err, "a", "b"), errB)
	checkIs(t, errB.With(err, "a", "b"), err)

	const errC = faults.Safe3[string, string, string]("a %s %s %s")
	check(t, errC, "a %s %s %s")
	check(t, errC.With(err, "a", "b", "c"), "^\\[github.com/fogfish/faults_test.TestSafe [0-9]+\\] a a b c: just error$")

	checkIs(t, errC, errC)
	checkIs(t, errC.With(err, "a", "b", "c"), errC)
	checkIs(t, errC.With(err, "a", "b", "c"), err)

	const errD = faults.Safe4[string, string, string, string]("a %s %s %s %s")
	check(t, errD, "a %s %s %s %s")
	check(t, errD.With(err, "a", "b", "c", "d"), "^\\[github.com/fogfish/faults_test.TestSafe [0-9]+\\] a a b c d: just error$")

	checkIs(t, errD, errD)
	checkIs(t, errD.With(err, "a", "b", "c", "d"), errD)
	checkIs(t, errD.With(err, "a", "b", "c", "d"), err)

	const errE = faults.Safe5[string, string, string, string, string]("a %s %s %s %s %s")
	check(t, errE, "a %s %s %s %s %s")
	check(t, errE.With(err, "a", "b", "c", "d", "e"), "^\\[github.com/fogfish/faults_test.TestSafe [0-9]+\\] a a b c d e: just error$")

	checkIs(t, errE, errE)
	checkIs(t, errE.With(err, "a", "b", "c", "d", "e"), errE)
	checkIs(t, errE.With(err, "a", "b", "c", "d", "e"), err)
}

// ------------------------------------------------------------------------------
//
// # Benchmark
//
// ------------------------------------------------------------------------------
var (
	err = fmt.Errorf("just error")
	exx error
	glo error
)

const (
	errFast = faults.Fast("error fast")
	errType = faults.Type("error type")
	errSafe = faults.Safe1[string]("error %s")
)

func failStdr() error { return fmt.Errorf("error type: %w", err) }
func failFast() error { return errFast.With(err) }
func failType() error { return errType.With(err) }
func failSafe() error { return errSafe.With(err, "safe") }

func BenchmarkStd(b *testing.B) {
	var err error

	for n := 0; n < b.N; n++ {
		exx = failStdr()
		if errors.Is(exx, err) {
			glo = err
		}
	}

	glo = err
}

func BenchmarkFast(b *testing.B) {
	var err error

	for n := 0; n < b.N; n++ {
		exx = failFast()
		if errors.Is(exx, errFast) {
			glo = err
		}
	}

	glo = err
}

func BenchmarkType(b *testing.B) {
	var err error

	for n := 0; n < b.N; n++ {
		exx = failType()
		if errors.Is(exx, errType) {
			glo = err
		}
	}

	glo = err
}

func BenchmarkSafe(b *testing.B) {
	var err error

	for n := 0; n < b.N; n++ {
		exx = failSafe()
		if errors.Is(exx, errSafe) {
			glo = err
		}
	}

	glo = err
}
