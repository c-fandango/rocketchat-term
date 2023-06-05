// Package utils for rocketchat-term
package utils

import (
	"math/rand"
	"strings"
	"unicode/utf8"
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
	for utf8.RuneCountInString(input) < n {
		input = padding + input
	}
	return input
}

func PadRight(input string, padding string, n int) string {
	for utf8.RuneCountInString(input) < n {
		input += padding
	}
	return input
}

func MergeStringMaps(first map[string]string, second map[string]string) map[string]string {
	for key, value := range second {
		first[key] = value
	}
	return first
}

func ReplaceEveryOther(input string, target string, replace string) string {
	inputSplt := strings.Split(input, target)
	for i := 1; i < len(inputSplt)-1; i += 2 {
		inputSplt[i] = replace + inputSplt[i] + target
	}
	return strings.Join(inputSplt, "")
}
