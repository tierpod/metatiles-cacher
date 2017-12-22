package lock

import (
	"fmt"
	"sort"
	"testing"
	"time"
)

func ExampleLock_Items() {
	l := New()
	l.Add("key1")
	l.Add("key2")

	items := l.Items()
	sort.Strings(items)
	fmt.Println(items)

	// Output:
	// [key1 key2]
}

func ExampleLock_Del() {
	l := New()
	l.Add("key1")
	l.Add("key2")
	l.Del("key1")
	fmt.Println(l.Items())

	// Output:
	// [key2]
}

func ExampleLock_HasKey() {
	l := New()
	l.Add("key1")
	fmt.Println(l.HasKey("key1"))
	fmt.Println(l.HasKey("key2"))

	// Output:
	// true
	// false
}

func ExampleLock_Wait() {
	l := New()

	// add item to locker
	l.Add("key1")
	// del item from locker after long-running job is done and send broadcast message to all waiters.
	defer l.Del("key1")

	// start second goroutine who will wait for broadcast message from Del
	go func() {
		if l.HasKey("key1") {
			fmt.Println("waiting for broadcast message from Del")
			l.Wait("key1", 1)
		}
	}()

	// emulate long-running work
	time.Sleep(30 * time.Millisecond)
	fmt.Println("done")

	// Output:
	// waiting for broadcast message from Del
	// done
}

func TestLockWaitTimeout(t *testing.T) {
	l := New()
	l.Add("key")
	defer l.Del("key")

	go func() {
		if l.HasKey("key") {
			err := l.Wait("key", 1)
			if err == nil {
				t.Errorf("expected %s, got nil error", ErrTimedOut)
			}
		}
	}()

	time.Sleep(2 * time.Second)
}
