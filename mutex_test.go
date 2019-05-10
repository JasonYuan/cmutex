package cmutex

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestMutex_Lock(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	m := NewMutex()
	assert.Nil(t, m.Lock(ctx))
	assert.True(t, m.Lock(ctx) != nil)
	m.Unlock()
	assert.Nil(t, m.Lock(context.TODO()))
}

func TestMutex_TryLock(t *testing.T) {
	m := NewMutex()
	assert.Nil(t, m.Lock(context.TODO()))
	assert.False(t, m.TryLock())
	m.Unlock()
	assert.True(t, m.TryLock())
}

var (
	rCount    = 100
	totalTime = time.Duration(rCount) * sleepTime
)

func TestLock2(t *testing.T) {
	start := time.Now()
	m := NewMutex()
	value := 0
	wg := sync.WaitGroup{}
	for i := 0; i < rCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = m.Lock(context.TODO())
			defer m.Unlock()
			assert.Equal(t, value, 0)
			value = 1
			time.Sleep(sleepTime)
			value = 0
		}()
	}
	wg.Wait()
	t.Log(time.Now().Sub(start), totalTime)
	assert.True(t, time.Now().Sub(start) >= totalTime)
}

var (
	cm = NewMutex()
	m  = sync.Mutex{}
)

func BenchmarkMutex_Lock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m.Lock()
		m.Unlock()
	}
}

func BenchmarkCMutex_Lock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cm.Lock(context.TODO())
		cm.Unlock()
	}
}

func BenchmarkMutex_Race(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		for j := 0; j < 50; j++ {
			wg.Add(1)
			go func() {
				m.Lock()
				m.Unlock()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkCMutex_Race(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		for j := 0; j < 50; j++ {
			wg.Add(1)
			go func() {
				_ = cm.Lock(context.TODO())
				cm.Unlock()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
