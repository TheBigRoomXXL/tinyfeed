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

func domain(item *gofeed.Item) string {
	url, err := url.Parse(item.Link)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return strings.TrimPrefix(url.Hostname(), "www.")
}

func publication(item *gofeed.Item) string {
	if item.PublishedParsed == nil {
		if item.Published != "" {
			return item.Published
		}
		return "Once upon a time"
	}
	return item.PublishedParsed.Format("2006-01-02")
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
