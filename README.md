# Promises for Go

## Overview

The Promises package provides an implementation of promises similar to JavaScript Promise for Go.  
This is a cleaner alternative to using channels, wait groups and goroutines for asynchronous operations.  
With this, you can write asynchronous code that looks like synchronous code.  
This is also compatible with the context package, so you can use them with context.Context.  
And this is written as `go-way`.

## Usage

### Importing the Package

Import the `promises` package in your Go code to use it:

```go
import "github.com/oneofthezombies/promises"
```

### Creating and Using Promises

You can create promises using the New function, passing an executor function that defines the asynchronous operation and how it should be resolved or rejected.

Here's a basic example:

```go
// promises_test.go
func TestAwait(t *testing.T) {
	ctx := context.Background()
	p := New(func(resolve Resolve[int], reject Reject) {
		resolve(1)
	})

	v, err := p.Await(ctx)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	if v != 1 {
		t.Errorf("expected value to be 1, got %d", v)
	}
}
```

## Additional Methods

Done: Get a channel that is closed when the promise is settled.  
Value: Get the value that the promise was fulfilled with.  
Reason (or Err for compatibility): Get the reason that the promise was rejected.  
IsFulfilled, IsRejected, IsSettled: Check the state of the promise.  

### All

You can use the All method to wait for multiple promises to be settled:

```go
// promises_test.go
func TestAllWithLoop(t *testing.T) {
	ctx := context.Background()
	var promises []*Promise[int]
	for i := 0; i < 100; i++ {
		i := i // https://golang.org/doc/faq#closures_and_goroutines
		p := New(func(resolve Resolve[int], reject Reject) {
			resolve(i)
		})

		promises = append(promises, p)
	}

	p := All(ctx, promises...)
	v, err := p.Await(ctx)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	if len(v) != 100 {
		t.Errorf("expected length to be 100, got %d", len(v))
	}

	for i, value := range v {
		if value != i {
			t.Errorf("expected value to be %d, got %d", i, value)
		}
	}
}
```

### AllSettled

You can use the AllSettled method to wait for multiple promises to be settled, regardless of whether they are resolved or rejected:

```go
// promises_test.go
func TestAllSettledWithLoop(t *testing.T) {
	ctx := context.Background()
	var promises []*Promise[int]
	for i := 0; i < 100; i++ {
		i := i // https://golang.org/doc/faq#closures_and_goroutines
		p := New(func(resolve Resolve[int], reject Reject) {
			resolve(i)
		})
		promises = append(promises, p)
	}

	p := AllSettled(ctx, promises...)

	v, err := p.Await(ctx)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	if len(v) != 100 {
		t.Errorf("expected length to be 100, got %d", len(v))
	}

	for i, result := range v {
		if result.Status != Fulfilled {
			t.Errorf("expected status to be Fulfilled, got %v", result.Status)
		}

		if result.Value != i {
			t.Errorf("expected value to be %d, got %d", i, result.Value)
		}
	}
}
```

## Contributing

If you would like to contribute to this project, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and ensure they pass the existing tests.
4. Add new tests if necessary.
5. Commit your changes with clear commit messages.
6. Push your branch to your fork.
7. Create a pull request to the original repository.

## License

This package is licensed under the MIT License. See the LICENSE file for details.
