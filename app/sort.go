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
	nameA, nameB := authorName(a), authorName(b)
	switch {
	case nameA == nameB:
		return compareByPublicationDate(a, b)
	case nameA == "":
		return 1
	case nameB == "":
		return -1
	}
	return strings.Compare(nameA, nameB)
}

// authorName returns the item's author name, or "" when there is no author.
func authorName(item Item) string {
	if item.Author == nil {
		return ""
	}
	return item.Author.Name
}
