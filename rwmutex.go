package cmutex

import (
	"context"
	"sync/atomic"
)

type rwMutex struct {
	wChan  chan struct{}
	rChan  chan struct{}
	rCount int32
}

type RWMutex interface {
	RLock(ctx context.Context) error
	Lock(ctx context.Context) error
	Unlock()
	RUnlock()
	TryRLock() bool
	TryLock() bool
}

func NewRWMutex() RWMutex {
	return &rwMutex{rChan: make(chan struct{}, 1), wChan: make(chan struct{}, 1), rCount: 0}
}

func (m *rwMutex) TryRLock() bool {
	select {
	case m.wChan <- struct{}{}:
	default:
		return false
	}
	atomic.AddInt32(&m.rCount, 1)
	<-m.wChan
	return true
}

func (m *rwMutex) RLock(ctx context.Context) error {
	select {
	case m.wChan <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	}
	atomic.AddInt32(&m.rCount, 1)
	<-m.wChan
	return nil
}

func (m *rwMutex) TryLock() bool {
	select {
	case m.wChan <- struct{}{}:
	default:
		return false
	}
	if atomic.LoadInt32(&m.rCount) <= 0 {
		return true
	}
	<-m.wChan
	return false
}

func (m *rwMutex) Lock(ctx context.Context) error {
	select {
	case m.wChan <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	}
	if atomic.LoadInt32(&m.rCount) <= 0 {
		return nil
	}
	for {
		select {
		case <-m.rChan:
			if atomic.LoadInt32(&m.rCount) <= 0 {
				return nil
			}
		case <-ctx.Done():
			<-m.wChan
			return ctx.Err()
		}
	}
}

func (m *rwMutex) Unlock() {
	select {
	case <-m.wChan:
	default:
	}
}

func (m *rwMutex) RUnlock() {
	after := atomic.AddInt32(&m.rCount, -1)
	if after < 0 {
		atomic.AddInt32(&m.rCount, 1)
	}
	if after <= 0 {
		select {
		case m.rChan <- struct{}{}:
		default:
		}
	}
}
