package production

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrBusy = errors.New("worker capacity exhausted")

type Dependency interface {
	Fetch(context.Context, string) (string, error)
}

type Service struct {
	dependency Dependency
	workers    chan struct{}
	stopping   chan struct{}
	once       sync.Once
}

func New(dependency Dependency, maxConcurrent int) *Service {
	return &Service{
		dependency: dependency,
		workers:    make(chan struct{}, maxConcurrent),
		stopping:   make(chan struct{}),
	}
}

// Lookup applies a service-owned deadline and refuses unbounded queued work.
func (s *Service) Lookup(ctx context.Context, id string) (string, error) {
	select {
	case <-s.stopping:
		return "", errors.New("service is shutting down")
	default:
	}
	select {
	case s.workers <- struct{}{}:
		defer func() { <-s.workers }()
	default:
		return "", ErrBusy
	}
	ctx, cancel := context.WithTimeout(ctx, 150*time.Millisecond)
	defer cancel()
	return s.dependency.Fetch(ctx, id)
}

func (s *Service) Shutdown() { s.once.Do(func() { close(s.stopping) }) }
