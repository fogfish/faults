# errors

The library `errors` provides type safe constructs to annotate Golang errors with the context and handle [opaque errors](https://tech.fog.fish/2022/07/05/assert-golang-errors-for-behavior.html#opaque-errors) without the boilerplate.

[![Version](https://img.shields.io/github/v/tag/fogfish/faults?label=version)](https://github.com/fogfish/faults/releases)
[![Documentation](https://pkg.go.dev/badge/github.com/fogfish/faults)](https://pkg.go.dev/github.com/fogfish/faults)
[![Build Status](https://github.com/fogfish/faults/workflows/test/badge.svg)](https://github.com/fogfish/faults/actions/)
[![Git Hub](https://img.shields.io/github/last-commit/fogfish/faults.svg)](https://github.com/fogfish/faults)
[![Coverage Status](https://coveralls.io/repos/github/fogfish/faults/badge.svg?branch=main)](https://coveralls.io/github/fogfish/faults?branch=main)


## Inspiration

The library is inspired by the post -- [Assert Golang Errors For Behavior: Everything You Need To Know Before Making Robust and Scalable Error Handling](https://github.com/fogfish/fogfish.github.io/blob/master/posts/2022/2022-07-05-assert-golang-errors-for-behavior.md) and its recommendation to deal with opaque errors -- failure does not have global catastrophic impacts but local functionality is impaired, execution of current control flow is terminated and incorrect results are returned.

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

In the context of a large application, `fmt.Errorf("[foo] dynamodb i/o failed: %w", err)` would be repeated a gazillion times for each service call to dynamodb, any refactoring of the error text or context structure becomes a tedious job - the type system would not help at all because `fmt.Errorf` is not type safe.

Usage of this library to define a type safe wrapping of errors is the better approach to annotate error context:

```go
// produces an error message
// [foo] dynamodb i/o failed: original error
const errDynamoIO = faults.Type("dynamodb i/o failed")

func foo() (*Bar, error) {
  val, err := db.dynamo.GetItem(ctx, req)
  if err != nil {
    return nil, errDynamoIO.With(err)
  }
  // continue happy path
}
```


// Type Safe + errors.Is




## Getting started

The latest version of the library is available at `main` branch of this repository. All development, including new features and bug fixes, take place on the `main` branch using forking and pull requests as described in contribution guidelines. The stable version is available via Golang modules.


### Error Type 

The `faults.Type` in the faults library provides a type-safe way to declare and manage error contexts in Go. This type allows you to define custom error types with specific messages, which can include placeholders for dynamic content. Here's how you can use faults.Type to create various types of errors.

```go
import "github.com/fogfish/faults"

const (
  // Create a basic error context with a static message
  ErrSomeA = faults.Type("something is failed")

  // Create an error context with a message that includes a placeholder for dynamic content
  ErrSomeB = faults.Type("something is failed %s")
)
```

Once errors are defined as constants, it used as

```go
func foo() error {
  return ErrSomeA
}

func bar() error {
  return ErrSomeB.With(err, "additional context")
}
```


### Type safe arguments

The `faults` library provides a way to create error contexts with type-safe arguments using `faults.Safe1`, `faults.Safe2`, up to `faults.Safe5`. These functions allow you to define error messages that include placeholders for dynamic content, ensuring that the arguments passed to the error message are type-checked at compile time.

```go
import "github.com/fogfish/faults"

const (
  // Create an error context with one type-safe argument
  errSomeD = faults.Safe1[int]("something %d is failed")

  // Create an error context with two type-safe arguments
  errSomeE = faults.Safe2[int, string]("something %d is failed %s")

  // Create an error context with three type-safe arguments
  errSomeF = faults.Safe3[int, string, float64]("something %d is failed %s with value %f")

  // Create an error context with four type-safe arguments
  errSomeG = faults.Safe4[int, string, float64, bool]("something %d is failed %s with value %f and flag %t")

  // Create an error context with five type-safe arguments
  errSomeH = faults.Safe5[int, string, float64, bool, string]("something %d is failed %s with value %f, flag %t and error %v")
)
```

Once you have defined the type-safe error contexts, you can use them in your functions to create detailed error messages with type-checked arguments.

```go
func foo() error {
  return errSomeD.With(err, 42)
}

func bar() error {
  return errSomeE.With(err, 42, "additional context")
}

func baz() error {
  return errSomeF.With(err, 42, "additional context", 3.14)
}

func qux() error {
  return errSomeG.With(err, 42, "additional context", 3.14, true)
}

func rug() error {
  return errSomeH.With(42, "additional context", 3.14, true, "all wrong")
}
```


### Matching errors

The faults library ensures that errors created using `faults.Type` and its variants are compatible with Go's standard `errors.Is` function. This allows you to match and compare errors effectively, making error handling more robust and consistent.

```go
import (
  "errors"
  "github.com/fogfish/faults"
)

const (
  errIO = faults.Type("service I/O has failed")
  errQuota = faults.Safe1[int]("quota %d has reached")
)
```

You can use `errors.Is` to check if an error matches a specific error type defined using the `faults` library. This is particularly useful for error handling and control flow in your application.

```go
func handleError(err error) {
  switch {
    case errors.Is(err, errQuota):
      // handle specific I/O error
    case errors.Is(err, errIO):
      // handle generic I/O 
    default:
      // handle other error
  }
}
```

By using `faults.Type` and its variants, you can create matchable errors that integrate seamlessly with Go's standard error handling mechanisms, making your error handling code more readable and maintainable.

### Errors for performance-critical apps 

The `faults` library provides a "fast" variant of error contexts using `faults.Fast`. This variant is designed for performance-critical applications where the overhead of using the `runtime` package to discover function context and inject it into the error is not acceptable. The `faults.Fast` type omits the usage of the runtime package internally, resulting in a significant performance improvement.

```go
import "github.com/fogfish/faults"

const (
  // Create a fast error context with a static message
  errFastA = faults.Fast("something is failed")

  // Create a fast error context with a message that includes a placeholder for dynamic content
  errFastB = faults.Fast("something is failed %s")
)
```

The `faults.Fast` type is particularly useful in highly loaded systems where the usage of the `runtime` package might cause about a 75% loss of the error path capacity. By omitting the `runtime` package, `faults.Fast` provides a more efficient way to handle errors without sacrificing the ability to include meaningful context in error messages.


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
git clone https://github.com/fogfish/faults
cd errors
go test
go test -run=^$ -bench=. -cpu 1 -benchtime=1s
```

### commit message

The commit message helps us to write a good release note, speed-up review process. The message should address two question what changed and why. The project follows the template defined by chapter [Contributing to a Project](http://git-scm.com/book/ch5-2.html) of Git book.

### bugs

If you experience any issues with the library, please let us know via [GitHub issues](https://github.com/fogfish/errors/issue). We appreciate detailed and accurate reports that help us to identity and replicate the issue. 


## License

[![See LICENSE](https://img.shields.io/github/license/fogfish/faults.svg?style=for-the-badge)](LICENSE)