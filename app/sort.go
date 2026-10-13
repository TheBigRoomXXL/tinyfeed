package tinyfeed

import "sort"

func sortItems(items []Item, orderBy string) []Item {
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
		if items[i].Author == nil {
			if items[j].Author == nil {
				return sortByPublicationDate(items)(i, j)
			}
			return false
		}
		if items[j].Author == nil {
			return true
		}
		if items[i].Author.Name == items[j].Author.Name {
			return sortByPublicationDate(items)(i, j)
		}
		return items[i].Author.Name < items[j].Author.Name
	}
}
