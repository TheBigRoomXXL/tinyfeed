package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"os"
	"strings"

	"github.com/mmcdole/gofeed"
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

func preview(item *gofeed.Item) string {
	if len(item.Description) > 0 {
		return item.Description
	}
	return item.Content
}

func domain(item *gofeed.Item) string {
	url, err := url.Parse(item.Link)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return strings.TrimPrefix(url.Hostname(), "www.")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Source: https://stackoverflow.com/questions/22892120/
func randStr(n int) string {
	const randomseed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = randomseed[rand.Int63()%int64(len(randomseed))]
	}
	return string(b)
}
