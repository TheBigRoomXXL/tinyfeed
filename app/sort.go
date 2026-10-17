package tinyfeed

import (
	"slices"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

func sortItems(items []Item, orderBy string) []Item {
	slices.SortStableFunc(items, itemComparator(orderBy))
	return items
}

// sortFeedItems sorts a single feed's items in place using the same order as
// sortItems, so that truncating them keeps the items that would be displayed.
func sortFeedItems(feed *gofeed.Feed, orderBy string) {
	cmp := itemComparator(orderBy)
	slices.SortStableFunc(feed.Items, func(a, b *gofeed.Item) int {
		return cmp(Item{a, feed.Title}, Item{b, feed.Title})
	})
}

func itemComparator(orderBy string) func(a, b Item) int {
	switch orderBy {
	case "publication-date":
		return compareChained(compareByPublicationDate, compareByUpdateDate, compareByAuthor, compareByFeedName, compareByLink)
	case "update-date":
		return compareChained(compareByUpdateDate, compareByPublicationDate, compareByAuthor, compareByFeedName, compareByLink)
	case "feed-name":
		return compareChained(compareByFeedName, compareByPublicationDate, compareByUpdateDate, compareByAuthor, compareByLink)
	case "author":
		return compareChained(compareByAuthor, compareByPublicationDate, compareByUpdateDate, compareByFeedName, compareByLink)
	default:
		panic("invalid orderBy")
	}
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

func compareString(a, b string) int {
	result := strings.Compare(strings.ToLower(a), strings.ToLower(b))
	if result != 0 {
		return result
	}
	return strings.Compare(a, b)
}

func compareByPublicationDate(a, b Item) int {
	return compareTimeDesc(a.PublishedParsed, b.PublishedParsed)
}

func compareByUpdateDate(a, b Item) int {
	return compareTimeDesc(a.UpdatedParsed, b.UpdatedParsed)
}

func compareByFeedName(a, b Item) int {
	return compareString(a.FeedName, b.FeedName)
}

func compareByLink(a, b Item) int {
	return compareString(a.Link, b.Link)
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
	return compareString(nameA, nameB)
}

// authorName returns the item's author name, or "" when there is no author.
func authorName(item Item) string {
	if item.Author == nil {
		return ""
	}
	return item.Author.Name
}
