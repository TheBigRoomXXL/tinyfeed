package tinyfeed

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"

	"github.com/mmcdole/gofeed"
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
	var less func(i, j int) bool

	switch orderBy {
	case "publication-date":
		less = sortByPublicationDate(items)
	case "update-date":
		less = sortByUpdateDate(items)
	case "feed-name":
		less = sortByFeedName(items)
	case "author":
		less = sortByAuthor(items)
	default:
		panic("invalid orderBy")
	}

	sort.SliceStable(items, less)
	return items
}

func sortByPublicationDate(items []Item) func(i, j int) bool {
	return func(i, j int) bool {
		if items[i].PublishedParsed == items[j].PublishedParsed {
			return sortByUpdateDate(items)(i, j)
		}
		if items[i].PublishedParsed == nil {
			return false
		}
		if items[j].PublishedParsed == nil {
			return true
		}
		return items[i].PublishedParsed.After(*items[j].PublishedParsed)
	}
}

func sortByUpdateDate(items []Item) func(i, j int) bool {
	return func(i, j int) bool {
		if items[i].UpdatedParsed == nil {
			return false
		}
		if items[j].UpdatedParsed == nil {
			return true
		}
		return items[i].UpdatedParsed.After(*items[j].UpdatedParsed)
	}
}

func sortByFeedName(items []Item) func(i, j int) bool {
	return func(i, j int) bool {
		if items[i].FeedName == items[j].FeedName {
			return sortByPublicationDate(items)(i, j)
		}
		return items[i].FeedName < items[j].FeedName
	}
}

func sortByAuthor(items []Item) func(i, j int) bool {
	return func(i, j int) bool {
		if items[i].Author.Name == items[j].Author.Name {
			return sortByPublicationDate(items)(i, j)
		}
		return items[i].Author.Name < items[j].Author.Name
	}
}

func opmlRSSVersion(feed gofeed.Feed) string {
	if feed.FeedType != "rss" {
		return feed.FeedVersion
	}
	switch feed.FeedVersion {
	case "1.0":
		return "RSS1"
	case "2.0":
		return "RSS2"
	default:
		return "RSS"
	}
}
