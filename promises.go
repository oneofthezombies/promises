package promises

import (
	"context"
	"errors"
	"sync"

	"github.com/oneofthezombies/option"
)

var (
	errReasonIsNil       = errors.New("reason must not be nil")
	errOptionValueIsNone = errors.New("option value must not be None")
)

type Resolve[T any] func(T)
type Reject func(error)
type Executor[T any] func(Resolve[T], Reject)

type OnFulfilled[T any] func(T)
type OnRejected func(error)
type OnFinally func()

// Reference: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise
type Promise[T any] struct {
	value  option.Option[T]
	reason error
	done   chan any
	mutex  sync.RWMutex
}

type Status int32

const (
	Fulfilled Status = iota
	Rejected
)

var statusStrings = [...]string{"fulfilled", "rejected"}

func (s Status) String() string {
	if s < Fulfilled || s > Rejected {
		return "unknown"
	}

	return statusStrings[s]
}

// New creates a new promise.
func New[T any](executor Executor[T]) *Promise[T] {
	p := &Promise[T]{
		value:  option.None[T](),
		reason: nil,
		done:   make(chan any),
	}

	resolve := func(value T) {
		p.mutex.Lock()
		defer p.mutex.Unlock()

		if p.isSettled() {
			return
		}

		defer close(p.done)
		p.value = option.Some(value)
	}

	reject := func(reason error) {
		if reason == nil {
			panic(errReasonIsNil)
		}

		p.mutex.Lock()
		defer p.mutex.Unlock()

		if p.isSettled() {
			return
		}

		defer close(p.done)
		p.reason = reason
	}

	go executor(resolve, reject)

	return p
}

func (p *Promise[T]) isFulfilled() bool {
	return p.value.Has()
}

func (p *Promise[T]) isRejected() bool {
	return p.reason != nil
}

func (p *Promise[T]) isSettled() bool {
	return p.isFulfilled() || p.isRejected()
}

// Then registers a callback that is called when the promise is fulfilled.
// If context is canceled, the callback is not called.
func (p *Promise[T]) Then(ctx context.Context, onFulfilled OnFulfilled[T]) *Promise[T] {
	select {
	case <-ctx.Done():
		return nil
	case <-p.done:
		break
	}

	p.mutex.RLock()
	v := p.value
	p.mutex.RUnlock()

	if !v.Has() {
		return p
	}

	onFulfilled(v.Value())
	return p
}

// Catch registers a callback that is called when the promise is rejected.
// If context is canceled, the callback is not called.
func (p *Promise[T]) Catch(ctx context.Context, onRejected OnRejected) *Promise[T] {
	select {
	case <-ctx.Done():
		return nil
	case <-p.done:
		break
	}

	p.mutex.RLock()
	r := p.reason
	p.mutex.RUnlock()

	if r == nil {
		return p
	}

	onRejected(r)
	return p
}

// Finally registers a callback that is called when the promise is settled.
// If context is canceled, the callback is not called.
func (p *Promise[T]) Finally(ctx context.Context, onFinally OnFinally) *Promise[T] {
	select {
	case <-ctx.Done():
		return nil
	case <-p.done:
		break
	}

	onFinally()
	return p
}

// Await blocks until the promise is settled and returns the value and reason or an error if the context is canceled.
func (p *Promise[T]) Await(ctx context.Context) (option.Option[T], error) {
	select {
	case <-ctx.Done():
		return option.None[T](), ctx.Err()
	case <-p.done:
		break
	}

	p.mutex.RLock()
	v := p.value
	r := p.reason
	p.mutex.RUnlock()

	return v, r
}

// Returns a channel that is closed when the promise is settled.
func (p *Promise[T]) Done() <-chan any {
	return p.done
}

// Get the value that the promise was fulfilled with.
// This method does not guarantee that the promise is settled.
// If you want to ensure that the promise is settled, use the Await() or Done() method before calling this method.
func (p *Promise[T]) Value() option.Option[T] {
	p.mutex.RLock()
	v := p.value
	p.mutex.RUnlock()

	return v
}

// Get the reason that the promise was rejected.
// This method does not guarantee that the promise is settled.
// If you want to ensure that the promise is settled, use the Await() or Done() method before calling this method.
func (p *Promise[T]) Reason() error {
	p.mutex.RLock()
	r := p.reason
	p.mutex.RUnlock()

	return r
}

// Get the reason that the promise was rejected.
// This method does not guarantee that the promise is settled.
// If you want to ensure that the promise is settled, use the Await() or Done() method before calling this method.
// This method is for compatibility with the Go standard library.
func (p *Promise[T]) Err() error {
	return p.Reason()
}

// Returns true if the promise is fulfilled.
// This method does not guarantee that the promise is settled.
// If you want to ensure that the promise is settled, use the Await() or Done() method before calling this method.
func (p *Promise[T]) IsFulfilled() bool {
	p.mutex.RLock()
	v := p.value
	p.mutex.RUnlock()

	return v.Has()
}

// Returns true if the promise is rejected.
// This method does not guarantee that the promise is settled.
// If you want to ensure that the promise is settled, use the Await() or Done() method before calling this method.
func (p *Promise[T]) IsRejected() bool {
	p.mutex.RLock()
	r := p.reason
	p.mutex.RUnlock()

	return r != nil
}

// Returns true if the promise is fulfilled or rejected.
// This method does not guarantee that the promise is settled.
// If you want to ensure that the promise is settled, use the Await() or Done() method before calling this method.
func (p *Promise[T]) IsSettled() bool {
	p.mutex.RLock()
	v := p.value
	r := p.reason
	p.mutex.RUnlock()

	return v.Has() || r != nil
}

// Reference: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise/all
func All[T any](ctx context.Context, promises ...*Promise[T]) *Promise[[]T] {
	p := New(func(resolve Resolve[[]T], reject Reject) {
		var wg sync.WaitGroup
		wg.Add(len(promises))

		results := make([]T, len(promises))
		for i, promise := range promises {
			go func(i int, promise *Promise[T]) {
				defer wg.Done()

				o, err := promise.Await(ctx)
				if err != nil {
					reject(err)
					return
				}

				if !o.Has() {
					reject(errOptionValueIsNone)
					return
				}

				results[i] = o.Value()
			}(i, promise)
		}

		wg.Wait()
		resolve(results)
	})

	return p
}

type SettledResult[T any] struct {
	Status Status
	Value  option.Option[T]
	Reason error
}

// Reference: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise/allSettled
func AllSettled[T any](ctx context.Context, promises ...*Promise[T]) *Promise[[]SettledResult[T]] {
	p := New(func(resolve Resolve[[]SettledResult[T]], reject Reject) {
		var wg sync.WaitGroup
		wg.Add(len(promises))

		results := make([]SettledResult[T], len(promises))
		for i, promise := range promises {
			go func(i int, promise *Promise[T]) {
				defer wg.Done()

				o, err := promise.Await(ctx)
				if err != nil {
					results[i] = SettledResult[T]{Status: Rejected, Reason: err}
					return
				}

				if !o.Has() {
					results[i] = SettledResult[T]{Status: Rejected, Reason: errOptionValueIsNone}
					return
				}

				results[i] = SettledResult[T]{Status: Fulfilled, Value: o}
			}(i, promise)
		}

		wg.Wait()
		resolve(results)
	})

	return p
}
