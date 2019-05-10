# cmutex

Mutex with context (timeout, cancel etc), also provide Trylock function similar to java lock.

Supporting mutex & rwmutex.

## Example

```go
package xxx

import (
	"context"
	"git.byted.org/ee/go/cmutex"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)
func TestMutex_Lock(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	m := cmutex.NewMutex()
	assert.Nil(t, m.Lock(ctx))
	assert.True(t, m.Lock(ctx) != nil)
	m.Unlock()
	assert.Nil(t, m.Lock(context.TODO()))
}

func TestMutex_TryLock(t *testing.T) {
	m := cmutex.NewMutex()
	assert.Nil(t, m.Lock(context.TODO()))
	assert.False(t, m.TryLock())
	m.Unlock()
	assert.True(t, m.TryLock())
}

func TestRwMutex_Lock(t *testing.T) {
	m := cmutex.NewRWMutex()
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
	m := cmutex.NewRWMutex()
	assert.Nil(t, m.RLock(context.TODO()))
	assert.True(t, m.TryRLock())
	assert.False(t, m.TryLock())
	m.RUnlock()
	m.RUnlock()
	assert.True(t, m.TryLock())
	assert.False(t, m.TryRLock())
}
```