package faults_test

import (
	"testing"

	"github.com/fogfish/faults"
)

func TestErrNotFound(t *testing.T) {
	const err = faults.ErrNotFound("key %s is not found")

	if faults.IsNotFound(err) {
		t.Fatalf("non initialized error cannot be used as NotFound")
	}

	exx := err.With(nil, "key")
	if !faults.IsNotFound(exx) {
		t.Fatalf("error has to be not found: %s", exx)
	}
}
