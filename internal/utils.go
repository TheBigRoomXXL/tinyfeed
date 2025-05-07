package tinyfeed

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/mmcdole/gofeed"
)

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
