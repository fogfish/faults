# errors

The library `errors` provides type safe constructs to annotate Golang errors with the context and handle [opaque errors](https://tech.fog.fish/2022/07/05/assert-golang-errors-for-behavior.html#opaque-errors) without the boilerplate.

[![Version](https://img.shields.io/github/v/tag/fogfish/errors?label=version)](https://github.com/fogfish/errors/releases)
[![Documentation](https://pkg.go.dev/badge/github.com/fogfish/errors)](https://pkg.go.dev/github.com/fogfish/errors)
[![Build Status](https://github.com/fogfish/errors/workflows/build/badge.svg)](https://github.com/fogfish/errors/actions/)
[![Git Hub](https://img.shields.io/github/last-commit/fogfish/errors.svg)](https://github.com/fogfish/errors)
[![Coverage Status](https://coveralls.io/repos/github/fogfish/errors/badge.svg?branch=main)](https://coveralls.io/github/fogfish/errors?branch=main)

## Inspiration

The library is inspired by the post -- [Assert Golang Errors For Behavior: Everything You Need To Know Before Making Robust and Scalable Error Handling](https://tech.fog.fish/2022/07/05/assert-golang-errors-for-behavior.html) and its recommendation to deal with opaque errors -- failure does not have global catastrophic impacts but local functionality is impaired, execution of current control flow is terminated and incorrect results are returned.

Opaque errors is a classical scenario of error handling in Golang, which is advertised by the majority of online publications. The code block knows that an error occurred. It does not have the ability to inspect error details, it only knows about the successful or unsuccessful completion of the called function. This error bubbles along the call stack until it is handler by the application.

```go
func foo() (*Bar, error) {
  val, err := db.dynamo.GetItem(ctx, req)
  if err != nil {
    return nil, err
  }
  // continue happy path
}
```

The debugging of opaque error handling becomes a difficult job because it violates **annotate errors with the context** recommendation. The Go Programming Language recommends inclusion of the context to the error path using `fmt.Errorf`.

```go
func foo() (*Bar, error) {
  val, err := db.dynamo.GetItem(ctx, req)
  if err != nil {
    return nil, fmt.Errorf("[foo] dynamodb i/o failed: %w", err)
  }
  // continue happy path
}
```

In the context of a large application, `fmt.Errorf("[foo] dynamodb i/o failed: %w", err)` would be repeated a gazillion times for each service call to dynamo, any refactoring of the error text or context structure becomes a tedious job - the type system would not help at all because `fmt.Errorf` is not type safe.

Usage of this library to define a type safe wrapping of errors is the better approach to annotate error context:

```go
// produces an error message
// [foo] dynamodb i/o failed: original error
const errDynamoIO = errors.Type("dynamodb i/o failed")

func foo() (*Bar, error) {
  val, err := db.dynamo.GetItem(ctx, req)
  if err != nil {
    return nil, errDynamoIO.New(err)
  }
  // continue happy path
}
```

## Getting started

The latest version of the library is available at `main` branch of this repository. All development, including new features and bug fixes, take place on the `main` branch using forking and pull requests as described in contribution guidelines. The stable version is available via Golang modules.

```go
import "github.com/fogfish/errors"

const (
  // create basic error context
  errSomeA = errors.Type("something is failed")
  // create error context with arguments
	errSomeB = errors.Type("something is failed %s")
  // create "fast" error context, would not annotate error with call stack
  errSomeC = errors.Fast("something is failed")
  // create error context with type safe arguments
	errSomeD = errors.Safe1[int]("something %d is failed")
  errSomeE = errors.Safe2[int, string]("something %d is failed %s")
)
```

### Gotchas 

The library uses the `runtime` package to discover function context and inject it into the error. If you are developing a highly loaded system, usage of `runtime` package might cause about 75% of the loss of the error path capacity. Therefore, the library support a "fast" variant of the type `errors.Fast`, which omits usage of `runtime` package internally.

## How To Contribute

The library is [MIT](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

The build and testing process requires [Go](https://golang.org) version 1.18 or later.

**build** and **test** library.

```bash
git clone https://github.com/fogfish/errors
cd errors
go test
go test -run=^$ -bench=. -cpu 1 -benchtime=1s
```

### commit message

The commit message helps us to write a good release note, speed-up review process. The message should address two question what changed and why. The project follows the template defined by chapter [Contributing to a Project](http://git-scm.com/book/ch5-2.html) of Git book.

### bugs

If you experience any issues with the library, please let us know via [GitHub issues](https://github.com/fogfish/errors/issue). We appreciate detailed and accurate reports that help us to identity and replicate the issue. 


## License

[![See LICENSE](https://img.shields.io/github/license/fogfish/errors.svg?style=for-the-badge)](LICENSE)