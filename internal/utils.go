package tinyfeed

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"sort"
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

func sortItems(items []Item) []Item {
	switch orderBy {
	case "publication-date":
		sort.SliceStable(items, func(i, j int) bool {
			if items[i].PublishedParsed == nil {
				return false
			}
			if items[j].PublishedParsed == nil {
				return true
			}
			return items[i].PublishedParsed.After(*items[j].PublishedParsed)
		})
	case "update-date":
		sort.SliceStable(items, func(i, j int) bool {
			if items[i].UpdatedParsed == nil {
				return false
			}
			if items[j].UpdatedParsed == nil {
				return true
			}
			return items[i].UpdatedParsed.After(*items[j].UpdatedParsed)
		})
	case "feed-name":
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].FeedName < items[j].FeedName
		})
	case "author":
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].Author.Name < items[j].Author.Name
		})
	}
	return items
}
