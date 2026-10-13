package tinyfeed

import (
	"slices"
	"strings"
	"time"
)

func sortItems(items []Item, orderBy string) []Item {
	var cmp func(a, b Item) int

	switch orderBy {
	case "publication-date":
		cmp = compareByPublicationDate
	case "update-date":
		cmp = compareByUpdateDate
	case "feed-name":
		cmp = compareByFeedName
	case "author":
		cmp = compareByAuthor
	default:
		panic("invalid orderBy")
	}

	slices.SortStableFunc(items, cmp)
	return items
}

// compareTimeDesc sorts newest first, nil last.
func compareTimeDesc(a, b *time.Time) int {
	switch {
	case a == nil && b == nil:
		return 0
	case a == nil:
		return 1
	case b == nil:
		return -1
	}
	return b.Compare(*a)
}

func compareByPublicationDate(a, b Item) int {
	if c := compareTimeDesc(a.PublishedParsed, b.PublishedParsed); c != 0 {
		return c
	}
	return compareByUpdateDate(a, b)
}

func compareByUpdateDate(a, b Item) int {
	return compareTimeDesc(a.UpdatedParsed, b.UpdatedParsed)
}

func compareByFeedName(a, b Item) int {
	if a.FeedName == b.FeedName {
		return compareByPublicationDate(a, b)
	}
	return strings.Compare(a.FeedName, b.FeedName)
}

func compareByAuthor(a, b Item) int {
	if a.Author == nil {
		if b.Author == nil {
			return compareByPublicationDate(a, b)
		}
		return 1
	}
	if b.Author == nil {
		return -1
	}
	if a.Author.Name == b.Author.Name {
		return compareByPublicationDate(a, b)
	}
	return strings.Compare(a.Author.Name, b.Author.Name)
}
