package queue

import (
	"fmt"
	"testing"
	"time"
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

func TestHasKey(t *testing.T) {
	q := NewUniq()
	q.Add("key")

	if !q.HasKey("key") {
		t.Errorf("HasKey: expected true, bot false")
	}
}

func TestWait(t *testing.T) {
	q := NewUniq()
	q.Add("key")
	err := q.Wait("key", 1)
	if err == nil {
		t.Errorf("Wait: expected error, got nil")
	}

	go func() {
		time.Sleep(1 * time.Second)
		q.Del("key")
	}()
	err = q.Wait("key", 2)
	if err != nil {
		t.Errorf("wait: expected nil, got error %v", err)
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
