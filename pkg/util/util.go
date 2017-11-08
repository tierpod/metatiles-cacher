// Package util contains useful functions.
package util

import (
	"crypto/md5"
	"fmt"
)

// MakeIntSlice makes slice of integers from min to max, not included max value.
func MakeIntSlice(min, max int) []int {
	r := make([]int, max-min)
	for i := range r {
		r[i] = min + i
	}
	return r
}

// DigestString returns md5sum of string
func DigestString(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
