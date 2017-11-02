package utils

import (
	"fmt"
	"testing"
)

func TestMakeIntRange(t *testing.T) {
	result := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := MakeIntRange(1, 10)

	if len(r) != len(result) {
		t.Errorf("MakeIntRange(1, 17): invalid result slice length (expected %v, got %v)", len(result), len(r))
	}

	for i := range r {
		if r[i] != result[i] {
			t.Errorf("MakeIntRange(1, 17): invalid slice item (expected: %v, got %v)", result[i], r[i])
		}
	}
}

func ExampleMakeIntRange() {
	r := MakeIntRange(1, 10)
	fmt.Printf("%v\n", r)

	// Output:
	// [1 2 3 4 5 6 7 8 9]
}

func TestMakeInludedIntRange(t *testing.T) {
	result := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	r := MakeIncludedIntRange(1, 10)

	if len(r) != len(result) {
		t.Errorf("MakeIntRange(1, 18): invalid result slice length (expected %v, got %v)", len(result), len(r))
	}

	for i := range r {
		if r[i] != result[i] {
			t.Errorf("MakeIntRange(1, 18): invalid slice item (expected: %v, got %v)", result[i], r[i])
		}
	}
}

func ExampleMakeIncludedIntRange() {
	r := MakeIncludedIntRange(1, 10)
	fmt.Printf("%v\n", r)

	// Output:
	// [1 2 3 4 5 6 7 8 9 10]
}

func ExampleDigestString() {
	d := DigestString("teststring")
	fmt.Printf("%v\n", d)

	// Output:
	// d67c5cbf5b01c9f91932e3b8def5e5f8
}
