package tinyfeed

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"
)

func domain(item Item) string {
	url, err := url.Parse(item.Link)
	if err != nil {
		log.Printf("WARNING: fail to parse url %s: %s\n", item.Link, err)
		return ""
	}
	return strings.TrimPrefix(url.Hostname(), "www.")
}

func publication(item Item) string {
	if item.PublishedParsed == nil {
		trimed := strings.TrimSpace(item.Published)
		if trimed == "" {
			return "Once upon a time"
		}
		return trimed
	}
	return item.PublishedParsed.Format("2006-01-02")
}

func generateNonce(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil { // Very unlikely
		log.Fatal(fmt.Errorf("failed to generate nonce: %w", err))
	}
	return base64.URLEncoding.EncodeToString(b)
}
