// Package utils contains utils functions
package utils

import "math/rand"

// MakeIncludedIntRange makes array of integers from min to max, included max value
func MakeIncludedIntRange(min, max int) []int {
	r := make([]int, max-min+1)
	for i := range r {
		r[i] = min + i
	}
	return r
}

// MakeIntRange makes array of integers from min to max, not included max value
func MakeIntRange(min, max int) []int {
	r := make([]int, max-min)
	for i := range r {
		r[i] = min + i
	}
	return r
}

// GetRandomItem choose random item from []string.
func GetRandomItem(items []string) string {
	var i int

	if len(items) > 1 {
		i = rand.Intn(len(items))
	}

	return items[i]
}
