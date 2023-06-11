// Package utils for rocketchat-term
package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"unicode/utf8"
)

func RandStr(n int) string {
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

func MaxInt(first int, second int) int {
	if first > second {
		return first
	}
	return second
}

func MinInt(first int, second int) int {
	if first < second {
		return first
	}
	return second
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

func MapperStr(input []string, f func(string) string) []string {
	output := make([]string, len(input), len(input))

	for i, item := range input {
		output[i] = f(item)
	}
	return output
}

func HexToRGB(hexCode string) (int, int, int, error) {
	if len(hexCode) != 6 {
		if len(hexCode) != 7 || hexCode[0] != '#' {
			return 0, 0, 0, fmt.Errorf("invalid hexcode: %s expecting e.g #ff123a", hexCode)
		}
		hexCode = hexCode[1:]
	}

	r, err := strconv.ParseInt(hexCode[:2], 16, 0)
	if err != nil {
		return 0, 0, 0, err
	}

	g, err := strconv.ParseInt(hexCode[2:4], 16, 0)
	if err != nil {
		return 0, 0, 0, err
	}

	b, err := strconv.ParseInt(hexCode[4:], 16, 0)
	if err != nil {
		return 0, 0, 0, err
	}

	return int(r), int(g), int(b), nil
}
