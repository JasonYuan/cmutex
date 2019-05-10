package cmutex

import (
	"context"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var cRWMutex = NewRWMutex()
var rWMutex = sync.RWMutex{}

func TestBasicRWMutex1(t *testing.T) {
	m := NewRWMutex()
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	assert.Nil(t, m.RLock(ctx))
	assert.Nil(t, m.RLock(ctx))
	assert.True(t, m.Lock(ctx) != nil)
	ctx, cancel = context.WithTimeout(context.TODO(), time.Second)
	m.RUnlock()
	m.RUnlock()
	assert.Nil(t, m.Lock(ctx))
	assert.True(t, m.RLock(ctx) != nil)
	assert.True(t, m.Lock(ctx) != nil)
	m.Unlock()
}

func TestRwMutex_TryLock(t *testing.T) {
	m := NewRWMutex()
	assert.Nil(t, m.RLock(context.TODO()))
	assert.True(t, m.TryRLock())
	assert.False(t, m.TryLock())
	m.RUnlock()
	m.RUnlock()
	assert.True(t, m.TryLock())
	assert.False(t, m.TryRLock())
}

var (
	readCount      = 100
	writeCount     = 100
	sleepTime      = 100 * time.Millisecond
	totalWriteTime = time.Duration(writeCount) * sleepTime
)

func TestRWMutex1(t *testing.T) {
	start := time.Now()
	m := NewRWMutex()
	value := 0
	wg := sync.WaitGroup{}
	for i := 0; i < writeCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.Nil(t, m.Lock(context.TODO()))
			defer func() {
				m.Unlock()
			}()
			assert.Equal(t, value, 0)
			value = 1
			time.Sleep(sleepTime)
			value = 0
		}()
	}
	for i := 0; i < readCount; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(rand.Int63n(int64(sleepTime))))
			defer wg.Done()
			assert.Nil(t, m.RLock(context.TODO()))
			//noinspection GoUnhandledErrorResult
			defer func() {
				m.RUnlock()
			}()
			assert.Equal(t, value, 0)
			time.Sleep(sleepTime)
			assert.Equal(t, value, 0)
		}()
	}
	wg.Wait()
	assert.True(t, time.Now().Sub(start) >= totalWriteTime)
}

func BenchmarkCRLock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cRWMutex.RLock(context.Background())
		cRWMutex.RUnlock()
	}
}

func BenchmarkRLock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rWMutex.RLock()
		rWMutex.RUnlock()
	}
}

func BenchmarkCLock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cRWMutex.Lock(context.Background())
		cRWMutex.Unlock()
	}
}

func BenchmarkLock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rWMutex.Lock()
		rWMutex.Unlock()
	}
}

func BenchmarkCRace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				_ = cRWMutex.RLock(context.Background())
				cRWMutex.RUnlock()
				wg.Done()
			}()
		}
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				_ = cRWMutex.Lock(context.Background())
				cRWMutex.Unlock()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkRace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				rWMutex.RLock()
				rWMutex.RUnlock()
				wg.Done()
			}()
		}
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				rWMutex.Lock()
				rWMutex.Unlock()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
