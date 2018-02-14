// Package util contains useful functions.
package util

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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

// Mimetype return mimetype based on file extension. Extension must be started with dot: `.png`
func Mimetype(ext string) (string, error) {
	switch ext {
	case ".png":
		return "image/png", nil
	case ".json", ".topojson", ".geojson":
		return "application/json", nil
	case ".mvt":
		return "application/vnd.mapbox-vector-tile", nil
	default:
		return "", fmt.Errorf("unknown mimetype for extension \"%v\"", ext)
	}
}

// MakeURL replaces {z}, {x}, {y} placeholders inside string `s` with values.
func MakeURL(s string, z, x, y int) string {
	url := strings.Replace(s, "{z}", strconv.Itoa(z), 1)
	url = strings.Replace(url, "{x}", strconv.Itoa(x), 1)
	url = strings.Replace(url, "{y}", strconv.Itoa(y), 1)

	return url
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString generates random string from english characters or digits.
func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
