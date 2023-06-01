// Package utils for rocketchat-term
package utils

import (
	"math/rand"
)

func RandID(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func PadLeft(input string, padding string, n int) string {
	for len(input) < n {
		input = padding + input
	}
	return input
}
