package production

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type slowDependency struct {
	started chan struct{}
	once    sync.Once
}

func (d *slowDependency) Fetch(ctx context.Context, _ string) (string, error) {
	d.once.Do(func() { close(d.started) })
	<-ctx.Done()
	return "", ctx.Err()
}

func TestDeadlineReachesDependency(t *testing.T) {
	service := New(&slowDependency{started: make(chan struct{})}, 1)
	started := time.Now()
	_, err := service.Lookup(context.Background(), "customer-1")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v", err)
	}
	if elapsed := time.Since(started); elapsed > 300*time.Millisecond {
		t.Fatalf("deadline took %s", elapsed)
	}
}

func TestWorkIsBoundedInsteadOfQueuingForever(t *testing.T) {
	dependency := &slowDependency{started: make(chan struct{})}
	service := New(dependency, 1)
	go func() { _, _ = service.Lookup(context.Background(), "first") }()
	<-dependency.started
	_, err := service.Lookup(context.Background(), "second")
	if !errors.Is(err, ErrBusy) {
		t.Fatalf("error = %v", err)
	}
}

func TestShutdownRejectsNewWork(t *testing.T) {
	service := New(&slowDependency{started: make(chan struct{})}, 1)
	service.Shutdown()
	_, err := service.Lookup(context.Background(), "customer-1")
	if err == nil {
		t.Fatal("expected shutdown rejection")
	}
}
