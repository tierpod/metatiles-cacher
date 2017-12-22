// Package lock provides wrapper around sync.Cond with map[string]interface{}.
package lock

import (
	"errors"
	"sync"
	"time"
)

var ErrTimedOut = errors.New("wait timed out")

// Lock is the basic struct with embedded *sync.Cond and map[string]interface{}
type Lock struct {
	cond  *sync.Cond
	items map[string]interface{}
}

// New creates new Lock.
func New() *Lock {
	m := sync.Mutex{}
	c := sync.NewCond(&m)
	i := make(map[string]interface{})

	l := &Lock{
		cond:  c,
		items: i,
	}

	return l
}

// Add adds key to Lock.
func (l *Lock) Add(key string) {
	l.cond.L.Lock()
	l.items[key] = nil
	l.cond.L.Unlock()
	// l.cond.Signal()
}

// Del dels key from Lock and send broadcast message to all waiter.
func (l *Lock) Del(key string) {
	l.cond.L.Lock()
	l.cond.Broadcast()
	delete(l.items, key)
	l.cond.L.Unlock()
}

// Wait waits for broadcast message. Use timeout in seconds for waiting. If timeout exceeded, return
// ErrTimedOut.
func (l *Lock) Wait(key string, timeout int) error {
	// create done channel for notify
	done := make(chan struct{})

	// run background waiter. close done channel if Wait() complete.
	go func() {
		l.cond.L.Lock()
		for hasKey(key, l.items) {
			l.cond.Wait()
		}
		l.cond.L.Unlock()
		close(done)
	}()

	// return ErrTimedOut if timeout exceeded
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		return ErrTimedOut
	// do nothing if Wait() complete successful
	case <-done:
	}

	return nil
}

// HasKey checks if key contains in Lock.
func (l *Lock) HasKey(key string) bool {
	l.cond.L.Lock()
	defer l.cond.L.Unlock()

	return hasKey(key, l.items)
}

// Items return all items stored in Lock.
func (l *Lock) Items() []string {
	n := []string{}

	l.cond.L.Lock()
	for k := range l.items {
		n = append(n, k)
	}
	l.cond.L.Unlock()

	return n
}

func hasKey(key string, items map[string]interface{}) bool {
	_, found := items[key]
	return found
}
