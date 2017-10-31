// Package queue implements queue with unique strings.
package queue

import "sync"

// Uniq contains mutex and map for storing strings.
type Uniq struct {
	mx sync.RWMutex
	m  map[string]bool
}

// NewUniq creates new Uniq queue.
func NewUniq() *Uniq {
	return &Uniq{
		m: make(map[string]bool),
	}
}

// Add adds key to map if item with this key does not exist (return true).
// Skip if key exist (return false).
func (q *Uniq) Add(key string) bool {
	q.mx.Lock()
	defer q.mx.Unlock()

	_, found := q.m[key]
	if !found {
		q.m[key] = true
		return true
	}

	return false
}

// Del deletes key from map.
func (q *Uniq) Del(key string) {
	q.mx.Lock()
	defer q.mx.Unlock()

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
	defer q.mx.RLocker()

	var result []string
	for k := range q.m {
		result = append(result, k)
	}

	return result
}
