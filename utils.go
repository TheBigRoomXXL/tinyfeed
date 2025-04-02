package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/mmcdole/gofeed"
)

func stdinToArgs() ([]string, error) {
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return readerToArgs(os.Stdin)
	}
	return []string{}, nil
}

func fileToArgs(filepath string) ([]string, error) {
	if filepath == "" {
		return []string{}, nil
	}
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening input file: %s", err)
	}
	return readerToArgs(file)
}

func readerToArgs(reader io.Reader) ([]string, error) {
	input, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading input: %s", err)
	}
	lines := strings.Split(string(input), "\n")

	argsSet := make(map[string]struct{})
	for _, s := range lines {
		lineTrimmed := strings.TrimSpace(s)
		if strings.HasPrefix(lineTrimmed, "#") || lineTrimmed == "" {
			continue
		}
		for _, field := range strings.Fields(lineTrimmed) {
			argsSet[field] = struct{}{}
		}
	}

	i := 0
	args := make([]string, len(argsSet))
	for k := range argsSet {
		args[i] = k
		i++
	}
	return args, nil
}

func domain(item *gofeed.Item) string {
	url, err := url.Parse(item.Link)
	if err != nil {
		log.Printf("WARNING: fail to parse url %s: %s\n", item.Link, err)
		return ""
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

func generateNonce(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil { // Very unlikely
		log.Fatal(fmt.Errorf("failed to generate nonce: %w", err))
	}
	return base64.URLEncoding.EncodeToString(b)
}
