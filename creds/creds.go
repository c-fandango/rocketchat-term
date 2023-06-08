package creds

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func GetUserInput(prompt string, secret bool) string {

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf(prompt)

	if secret {
		inputSecret, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			panic("error in reading input")
		}
		return strings.TrimSpace(string(inputSecret))
	}

	userInput, err := reader.ReadString('\n')

	if err != nil {
		panic("error in reading input")
	}

	return strings.TrimSpace(userInput)
}

func WriteCache(path string, cache []byte) error {
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	err := os.WriteFile(path, cache, 0644)

	if err != nil {
		return fmt.Errorf("failed write cache to: %w", err)
	}

	return nil
}

func ReadCache(path string) ([]byte, error) {
	b, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("no cache found at location %s", path)
	}

	return b, nil
}

func ClearCache(path string) {
	os.Remove(path)
}
