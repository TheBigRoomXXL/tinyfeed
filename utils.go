package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func stdinToArgs() ([]string, error) {
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("error parsing stdin: %s", err)
		}
		return strings.Fields(string(input)), nil
	}
	return []string{}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func truncstr(input string, length int) string {
	asRunes := []rune(input)

	if length > len(asRunes) {
		return input
	} else {
		return string(asRunes[0:length]) + "..."
	}
}
