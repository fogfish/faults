//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/errors
//

package main

import (
	"fmt"
	"os"

	"github.com/fogfish/faults"
)

const (
	errDoSomething     = faults.Type("unable to do something")
	errDoSomethingElse = faults.Type("unable to do something else (e.g. %s)")
	errDoAttempt       = faults.Safe1[int]("attempt %d is failed")
)

func doSomething() error {
	_, err := os.Open("some-file-which-do-exist")
	if err != nil {
		return errDoSomething.New(err)
	}

	return nil
}

func doSomethingElse() error {
	if err := doSomething(); err != nil {
		return errDoSomethingElse.New(err, "nested call")
	}

	return nil
}

func doAttempt(n int) error {
	if err := doSomethingElse(); err != nil {
		return errDoAttempt.New(err, n)
	}

	return nil

}

func main() {
	fmt.Println(doAttempt(10))
}
