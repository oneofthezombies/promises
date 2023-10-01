# Promises for Go

## Overview

The Go Promises package provides an implementation of promises similar to JavaScript promises. Promises are used for handling asynchronous operations and managing callbacks.

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
ctx := context.Background()
promise := promises.New[int](func(resolve promises.Resolve[int], reject promises.Reject) {
    // Perform asynchronous operation
    result := 42 // Simulated result
    
    // Resolve the promise with the result
    resolve(result)
})

// Attach callbacks using Then and Catch
promise.Then(ctx, func(value int) {
    fmt.Printf("Promise resolved with value: %d\n", value)
}).Catch(ctx, func(err error) {
    fmt.Printf("Promise rejected with error: %s\n", err.Error())
})
```

### Handling Errors

You can use Catch to handle errors in promises:

```go
ctx := context.Background()
promise := promises.New[int](func(resolve promises.Resolve[int], reject promises.Reject) {
    // Simulate an error
    err := errors.New("An error occurred")
    
    // Reject the promise with the error
    reject(err)
})

promise.Catch(ctx, func(err error) {
    fmt.Printf("Promise rejected with error: %s\n", err.Error())
})
```

### Await and Synchronization

You can use the Await method to block until the promise is settled:

```go
ctx := context.Background()
result, err := promise.Await(ctx)
if err != nil {
    fmt.Printf("Promise rejected with error: %s\n", err.Error())
} else {
    fmt.Printf("Promise resolved with value: %d\n", result.Value())
}
```

## Additional Methods

Finally: Register a callback that is called when the promise is settled.  
Done: Get a channel that is closed when the promise is settled.  
Value: Get the value that the promise was fulfilled with.  
Reason (or Err for compatibility): Get the reason that the promise was rejected.  
IsFulfilled, IsRejected, IsSettled: Check the state of the promise.  

### All

You can use the All method to wait for multiple promises to be settled:

```go
ctx := context.Background()
p1 := New(func(resolve Resolve[int], reject Reject) {
    resolve(1)
})
p2 := New(func(resolve Resolve[int], reject Reject) {
    resolve(2)
})
p3 := New(func(resolve Resolve[int], reject Reject) {
    resolve(3)
})

p := All(ctx, p1, p2, p3)
result, err := p.Await(ctx)
// handle result and error. please, see promises_test.go TestAll()
```

### AllSettled

You can use the AllSettled method to wait for multiple promises to be settled, regardless of whether they are resolved or rejected:

```go
ctx := context.Background()
p1 := New(func(resolve Resolve[int], reject Reject) {
    resolve(1)
})
p2 := New(func(resolve Resolve[int], reject Reject) {
    reject(errors.New("something went wrong"))
})
p3 := New(func(resolve Resolve[int], reject Reject) {
    resolve(3)
})

p := All(ctx, p1, p2, p3)

result, err := p.Await(ctx)
// handle result and error. please, see promises_test.go TestAllSettled()
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
