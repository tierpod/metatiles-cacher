package queue

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	q := NewUniq()

	added := q.Add("key")
	if added == false {
		t.Errorf("Add: expected true, got false")
	}

	added = q.Add("key")
	if added == true {
		t.Errorf("Add item with exist key: expected false, got true")
	}
}

func TestDel(t *testing.T) {
	q := NewUniq()

	q.Add("key")
	q.Del("key")

	if q.Len() != 0 {
		t.Errorf("Del: expected length 0, got %v", q.Len())
	}
}

func ExampleUniq_Items() {
	q := NewUniq()
	q.Add("key1")
	q.Add("key2")
	fmt.Printf("q.Length: %v\n", q.Items())

	// Output:
	// q.Length: [key1 key2]
}
