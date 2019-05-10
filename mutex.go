package cmutex

import "context"

type mutex struct {
	c chan struct{}
}

type Mutex interface {
	TryLock() bool
	Lock(ctx context.Context) error
	Unlock()
}

func NewMutex() Mutex {
	return &mutex{c: make(chan struct{}, 1)}
}

func (m *mutex) TryLock() bool {
	select {
	case m.c <- struct{}{}:
		return true
	default:
		return false
	}
}

func (m *mutex) Lock(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case m.c <- struct{}{}:
		return nil
	}
}

func (m *mutex) Unlock() {
	select {
	case <-m.c:
	default:
	}
}
