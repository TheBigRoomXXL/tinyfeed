package tinyfeed

import (
	"slices"
	"strings"
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

func compareByPublicationDate(a, b Item) int {
	if a.PublishedParsed == b.PublishedParsed {
		return compareByUpdateDate(a, b)
	}
	if a.PublishedParsed == nil {
		return 1
	}
	if b.PublishedParsed == nil {
		return -1
	}
	return b.PublishedParsed.Compare(*a.PublishedParsed)
}

func compareByUpdateDate(a, b Item) int {
	if a.UpdatedParsed == nil && b.UpdatedParsed == nil {
		return 0
	}
	if a.UpdatedParsed == nil {
		return 1
	}
	if b.UpdatedParsed == nil {
		return -1
	}
	return b.UpdatedParsed.Compare(*a.UpdatedParsed)
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
