package promises_test

import (
	"context"
	"errors"
	"testing"
	"time"

	. "github.com/oneofthezombies/promises"
)

func TestNew(t *testing.T) {
	p := New(func(resolve Resolve[int], reject Reject) {
		resolve(1)
	})

	if p == nil {
		t.Errorf("expected promise to be non-nil")
	}
}

func TestThen(t *testing.T) {
	ctx := context.Background()
	p := New(func(resolve Resolve[int], reject Reject) {
		resolve(1)
	})

	p.Then(ctx, func(value int) {
		if value != 1 {
			t.Errorf("expected value to be 1, got %d", value)
		}
	}).Catch(ctx, func(reason error) {
		t.Errorf("expected reason to be nil, got %v", reason)
	})
}

func TestCatch(t *testing.T) {
	ctx := context.Background()
	p := New(func(resolve Resolve[int], reject Reject) {
		reject(errors.New("something went wrong"))
	})

	p.Then(ctx, func(value int) {
		t.Errorf("expected value to be nil, got %d", value)
	}).Catch(ctx, func(reason error) {
		t.Logf("reason: %v", reason)
	})
}

func TestFinally(t *testing.T) {
	ctx := context.Background()
	p := New(func(resolve Resolve[int], reject Reject) {
		resolve(1)
	})

	p.Finally(ctx, func() {
		t.Logf("finally")
	})
}

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

func TestAwaitMultipleResolve(t *testing.T) {
	ctx := context.Background()
	p := New(func(resolve Resolve[int], reject Reject) {
		resolve(1)
		resolve(2)
	})

	v, err := p.Await(ctx)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	if v != 1 {
		t.Errorf("expected value to be 1, got %d", v)
	}
}

func TestAwaitWithGo(t *testing.T) {
	ctx := context.Background()
	p := New(func(resolve Resolve[int], reject Reject) {
		go func() {
			time.Sleep(1 * time.Second)
			resolve(1)
		}()
	})

	v, err := p.Await(ctx)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	if v != 1 {
		t.Errorf("expected value to be 1, got %d", v)
	}
}

func TestIsSettled(t *testing.T) {
	p := New(func(resolve Resolve[int], reject Reject) {
		resolve(1)
	})

	<-p.Done()
	if !p.IsSettled() {
		t.Errorf("expected promise to be settled")
	}
}

func TestAwaitCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	p := New(func(resolve Resolve[int], reject Reject) {
		time.Sleep(3 * time.Second)
		resolve(1)
	})

	cancel()
	_, err := p.Await(ctx)
	if err == nil {
		t.Errorf("expected error to be non-nil")
	}
}

func TestAwaitString(t *testing.T) {
	ctx := context.Background()
	p := New(func(resolve Resolve[string], reject Reject) {
		resolve("hello")
	})

	v, err := p.Await(ctx)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	if v != "hello" {
		t.Errorf("expected value to be hello, got %s", v)
	}
}

func TestAwaitError(t *testing.T) {
	ctx := context.Background()
	p := New(func(resolve Resolve[string], reject Reject) {
		reject(errors.New("something went wrong"))
	})

	_, err := p.Await(ctx)
	if err == nil {
		t.Errorf("expected error to be non-nil")
	}
}

func TestAwaitStruct(t *testing.T) {
	ctx := context.Background()
	p := New(func(resolve Resolve[struct{ Name string }], reject Reject) {
		resolve(struct{ Name string }{Name: "hello"})
	})

	v, err := p.Await(ctx)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	if v.Name != "hello" {
		t.Errorf("expected value to be hello, got %s", v.Name)
	}
}

func TestAll(t *testing.T) {
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

	v, err := p.Await(ctx)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	if len(v) != 3 {
		t.Errorf("expected length to be 3, got %d", len(v))
	}

	if v[0] != 1 {
		t.Errorf("expected value to be 1, got %d", v[0])
	}

	if v[1] != 2 {
		t.Errorf("expected value to be 2, got %d", v[1])
	}

	if v[2] != 3 {
		t.Errorf("expected value to be 3, got %d", v[2])
	}
}

func TestAllWithRejected(t *testing.T) {
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

	_, err := p.Await(ctx)
	if err == nil {
		t.Errorf("expected error to be non-nil")
	}
}

func TestAllSettled(t *testing.T) {
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

	p := AllSettled(ctx, p1, p2, p3)

	v, err := p.Await(ctx)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	if len(v) != 3 {
		t.Errorf("expected length to be 3, got %d", len(v))
	}

	if v[0].Status != Fulfilled {
		t.Errorf("expected status to be Fulfilled, got %v", v[0].Status)
	}

	if v[0].Value != 1 {
		t.Errorf("expected value to be 1, got %d", v[0].Value)
	}

	if v[1].Status != Rejected {
		t.Errorf("expected status to be Rejected, got %v", v[1].Status)
	}

	if v[1].Reason == nil {
		t.Errorf("expected reason to be non-nil")
	}

	if v[2].Status != Fulfilled {
		t.Errorf("expected status to be Fulfilled, got %v", v[2].Status)
	}

	if v[2].Value != 3 {
		t.Errorf("expected value to be 3, got %d", v[2].Value)
	}
}

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
