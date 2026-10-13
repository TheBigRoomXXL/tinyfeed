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
		cmp = compareChained(compareByPublicationDate, compareByUpdateDate, compareByAuthor, compareByFeedName)
	case "update-date":
		cmp = compareChained(compareByUpdateDate, compareByPublicationDate, compareByAuthor, compareByFeedName)
	case "feed-name":
		cmp = compareChained(compareByFeedName, compareByPublicationDate, compareByUpdateDate, compareByAuthor)
	case "author":
		cmp = compareChained(compareByAuthor, compareByPublicationDate, compareByUpdateDate, compareByFeedName)
	default:
		panic("invalid orderBy")
	}

	slices.SortStableFunc(items, cmp)
	return items
}

// compareChained returns a comparison that applies each cmp in order and
// returns the first non-zero result.
func compareChained(cmps ...func(a, b Item) int) func(a, b Item) int {
	return func(a, b Item) int {
		for _, cmp := range cmps {
			if c := cmp(a, b); c != 0 {
				return c
			}
		}
		return 0
	}
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
	return compareTimeDesc(a.PublishedParsed, b.PublishedParsed)
}

func compareByUpdateDate(a, b Item) int {
	return compareTimeDesc(a.UpdatedParsed, b.UpdatedParsed)
}

func compareByFeedName(a, b Item) int {
	return strings.Compare(a.FeedName, b.FeedName)
}

func compareByAuthor(a, b Item) int {
	nameA, nameB := authorName(a), authorName(b)
	switch {
	case nameA == nameB:
		return 0
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
