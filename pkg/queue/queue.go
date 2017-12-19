// Package queue implements queue with unique strings.
package queue

import (
	"errors"
	"sync"
	"time"
)

// Uniq contains mutex and map for storing strings. Map contains string as `key` and `chan bool` as
// value.
type Uniq struct {
	mx sync.RWMutex
	m  map[string]chan bool
}

// NewUniq creates new Uniq queue.
func NewUniq() *Uniq {
	return &Uniq{
		m: make(map[string]chan bool),
	}
}

// Add adds item with given key to map if item with this key does not exist (return true).
// Skip if key exist (return false).
func (q *Uniq) Add(key string) {
	q.mx.Lock()
	defer q.mx.Unlock()

	_, found := q.m[key]
	if found {
		return
	}

	done := make(chan bool)
	q.m[key] = done
}

// Del deletes key from map.
func (q *Uniq) Del(key string) {
	q.mx.Lock()
	defer q.mx.Unlock()

	close(q.m[key])
	delete(q.m, key)
}

// Len returns map length.
func (q *Uniq) Len() int {
	q.mx.RLock()
	defer q.mx.RUnlock()
	return len(q.m)
}

// Items returns slice of keys storing in map.
func (q *Uniq) Items() []string {
	q.mx.RLock()
	defer q.mx.RUnlock()

	var result []string
	for k := range q.m {
		result = append(result, k)
	}

	return result
}

// HasKey checks if queue has item with key.
func (q *Uniq) HasKey(key string) bool {
	q.mx.RLock()
	defer q.mx.RUnlock()

	_, found := q.m[key]
	return found
}

// ErrWaitTimeout is the error if Wait timeout achieved.
var ErrWaitTimeout = errors.New("wait timeout")

// Wait waits until key was deleted from queue or timeout appears.
func (q *Uniq) Wait(key string, timeout int) error {
	if q.HasKey(key) {
		select {
		case <-q.m[key]:
			// fmt.Printf("DONE CHAN CLOSED FOR KEY: %v\n", key)
		case <-time.After(time.Second * time.Duration(timeout)):
			return ErrWaitTimeout
		}
	}

	return nil
}
