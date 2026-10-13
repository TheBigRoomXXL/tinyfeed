package tinyfeed

import (
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
)

func TestSortItems(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	lastWeek := now.AddDate(0, 0, -7)

	itemA := Item{
		Item: &gofeed.Item{
			Title:           "Alpha",
			PublishedParsed: &yesterday,
			UpdatedParsed:   &now,
			Author:          &gofeed.Person{Name: "Charlie"},
		},
		FeedName: "Z-Feed",
	}

	itemB := Item{
		Item: &gofeed.Item{
			Title:           "Beta",
			PublishedParsed: &now,
			UpdatedParsed:   &lastWeek,
			Author:          &gofeed.Person{Name: "Alice"},
		},
		FeedName: "A-Feed",
	}

	itemC := Item{
		Item: &gofeed.Item{
			Title:           "Gamma",
			PublishedParsed: &lastWeek,
			UpdatedParsed:   &yesterday,
			Author:          &gofeed.Person{Name: "Bob"},
		},
		FeedName: "M-Feed",
	}

	tests := []struct {
		label    string
		orderBy  string
		input    []Item
		expected []string // We'll check the Titles to verify order
	}{
		{
			label:    "Sort by publication-date (descending)",
			orderBy:  "publication-date",
			input:    []Item{itemC, itemA, itemB},
			expected: []string{"A-Feed", "Z-Feed", "M-Feed"},
		},
		{
			label:    "Sort by update-date (descending)",
			orderBy:  "update-date",
			input:    []Item{itemC, itemA, itemB},
			expected: []string{"Z-Feed", "M-Feed", "A-Feed"},
		},
		{
			label:    "Sort by feed-name (ascending)",
			orderBy:  "feed-name",
			input:    []Item{itemA, itemB, itemC},
			expected: []string{"A-Feed", "M-Feed", "Z-Feed"},
		},
		{
			label:    "Sort by author (ascending)",
			orderBy:  "author",
			input:    []Item{itemA, itemB, itemC},
			expected: []string{"A-Feed", "M-Feed", "Z-Feed"},
		},
	}

	for _, test := range tests {
		t.Run(test.label, func(t *testing.T) {
			got := sortItems(test.input, test.orderBy)

			for i, item := range got {
				if item.FeedName != test.expected[i] {
					t.Errorf("%s: at index %d, expected %s, got %s", test.label, i, test.expected[i], item.FeedName)
				}
			}
		})
	}
}

func TestSortItemsWithNilDates(t *testing.T) {
	now := time.Now()
	itemNil := Item{Item: &gofeed.Item{Title: "NilDate", PublishedParsed: nil}}
	itemValid := Item{Item: &gofeed.Item{Title: "ValidDate", PublishedParsed: &now}}

	items := []Item{itemNil, itemValid}
	sorted := sortItems(items, "publication-date")

	if sorted[0].Title != "ValidDate" {
		t.Errorf("Expected ValidDate to come before NilDate in descending sort")
	}
}

func TestSortItemsWithMissingAuthors(t *testing.T) {
	newer := time.Date(2026, time.January, 2, 0, 0, 0, 0, time.UTC)
	older := newer.AddDate(0, 0, -1)
	named := func(title, author string) Item {
		return Item{Item: &gofeed.Item{Title: title, Author: &gofeed.Person{Name: author}}}
	}
	missing := func(title string, published *time.Time) Item {
		return Item{Item: &gofeed.Item{Title: title, PublishedParsed: published}}
	}
	tests := []struct {
		name  string
		items []Item
		want  []string
	}{
		{
			name:  "missing before named",
			items: []Item{missing("Missing", nil), named("Alice", "Alice")},
			want:  []string{"Alice", "Missing"},
		},
		{
			name:  "named before missing",
			items: []Item{named("Alice", "Alice"), missing("Missing", nil)},
			want:  []string{"Alice", "Missing"},
		},
		{
			name:  "missing authors use publication date",
			items: []Item{missing("Older", &older), missing("Newer", &newer)},
			want:  []string{"Newer", "Older"},
		},
		{
			name:  "missing authors and dates stay stable",
			items: []Item{missing("First", nil), missing("Second", nil)},
			want:  []string{"First", "Second"},
		},
		{
			name:  "missing authors with undated entries",
			items: []Item{missing("Undated", nil), missing("Dated", &older)},
			want:  []string{"Dated", "Undated"},
		},
		{
			name: "mixed authors",
			items: []Item{
				missing("Older", &older), named("Bob", "Bob"),
				missing("Newer", &newer), named("Alice", "Alice"),
			},
			want: []string{"Alice", "Bob", "Newer", "Older"},
		},
		{
			name:  "empty author name is treated as missing",
			items: []Item{missing("Missing", nil), named("Alice", "Alice"), named("Empty", "")},
			want:  []string{"Alice", "Missing", "Empty"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := sortItems(test.items, "author")
			for i, want := range test.want {
				if got[i].Title != want {
					t.Errorf("item %d = %q, want %q", i, got[i].Title, want)
				}
			}
		})
	}
}

// Each case ties on every criterion before the one under test, and the later
// criteria disagree with the expected order, so only that criterion can
// decide it.
func TestSortItemsFallbackChain(t *testing.T) {
	newer := time.Date(2026, time.January, 2, 0, 0, 0, 0, time.UTC)
	older := newer.AddDate(0, 0, -1)
	item := func(title, feed, author string, published, updated *time.Time) Item {
		return Item{
			Item: &gofeed.Item{
				Title:           title,
				Author:          &gofeed.Person{Name: author},
				PublishedParsed: published,
				UpdatedParsed:   updated,
			},
			FeedName: feed,
		}
	}
	tests := []struct {
		orderBy string
		name    string
		first   Item
		second  Item
	}{
		{
			"publication-date", "then update date",
			item("First", "B", "Bob", &newer, &newer),
			item("Second", "A", "Alice", &newer, &older),
		},
		{
			"publication-date", "then author",
			item("First", "B", "Alice", &newer, &newer),
			item("Second", "A", "Bob", &newer, &newer),
		},
		{
			"publication-date", "then feed name",
			item("First", "A", "Alice", &newer, &newer),
			item("Second", "B", "Alice", &newer, &newer),
		},
		{
			"update-date", "then publication date",
			item("First", "B", "Bob", &newer, &newer),
			item("Second", "A", "Alice", &older, &newer),
		},
		{
			"update-date", "then author",
			item("First", "B", "Alice", &newer, &newer),
			item("Second", "A", "Bob", &newer, &newer),
		},
		{
			"update-date", "then feed name",
			item("First", "A", "Alice", &newer, &newer),
			item("Second", "B", "Alice", &newer, &newer),
		},
		{
			"feed-name", "then publication date",
			item("First", "A", "Bob", &newer, &older),
			item("Second", "A", "Alice", &older, &newer),
		},
		{
			"feed-name", "then update date",
			item("First", "A", "Bob", &newer, &newer),
			item("Second", "A", "Alice", &newer, &older),
		},
		{
			"feed-name", "then author",
			item("First", "A", "Alice", &newer, &newer),
			item("Second", "A", "Bob", &newer, &newer),
		},
		{
			"author", "then publication date",
			item("First", "B", "Alice", &newer, &older),
			item("Second", "A", "Alice", &older, &newer),
		},
		{
			"author", "then update date",
			item("First", "B", "Alice", &newer, &newer),
			item("Second", "A", "Alice", &newer, &older),
		},
		{
			"author", "then feed name",
			item("First", "A", "Alice", &newer, &newer),
			item("Second", "B", "Alice", &newer, &newer),
		},
	}
	for _, test := range tests {
		t.Run(test.orderBy+" "+test.name, func(t *testing.T) {
			got := sortItems([]Item{test.second, test.first}, test.orderBy)
			if got[0].Title != "First" || got[1].Title != "Second" {
				t.Errorf("got [%q, %q], want [First, Second]", got[0].Title, got[1].Title)
			}
		})
	}
}

func TestSortItemsByFeedNameAlphabetical(t *testing.T) {
	feeds := []string{"b", "Feed10", "Feed2", "", "Feed", "B", "a"}
	items := make([]Item, len(feeds))
	for i, feed := range feeds {
		items[i] = Item{Item: &gofeed.Item{}, FeedName: feed}
	}
	want := []string{"", "B", "Feed", "Feed10", "Feed2", "a", "b"}
	got := sortItems(items, "feed-name")
	for i, w := range want {
		if got[i].FeedName != w {
			t.Errorf("item %d = %q, want %q", i, got[i].FeedName, w)
		}
	}
}

// Regression for 38b93c9: equal dates held by distinct pointers must be
// treated as equal and fall back to the next criterion.
func TestSortItemsCompareDateValuesNotPointers(t *testing.T) {
	newer := time.Date(2026, time.January, 2, 0, 0, 0, 0, time.UTC)
	older := newer.AddDate(0, 0, -1)
	publishedA := newer
	publishedB := newer // same value, different pointer

	itemA := Item{Item: &gofeed.Item{Title: "UpdatedOlder", PublishedParsed: &publishedA, UpdatedParsed: &older}}
	itemB := Item{Item: &gofeed.Item{Title: "UpdatedNewer", PublishedParsed: &publishedB, UpdatedParsed: &newer}}

	tests := []struct {
		orderBy string
		items   []Item
	}{
		{"publication-date", []Item{itemA, itemB}},
		{"feed-name", []Item{itemA, itemB}},
		{"author", []Item{itemA, itemB}},
	}
	for _, test := range tests {
		t.Run(test.orderBy, func(t *testing.T) {
			got := sortItems(test.items, test.orderBy)
			if got[0].Title != "UpdatedNewer" || got[1].Title != "UpdatedOlder" {
				t.Errorf("got [%q, %q], want [UpdatedNewer, UpdatedOlder]", got[0].Title, got[1].Title)
			}
		})
	}
}

// Regression for 489566a: an empty author name is equivalent to a missing
// author, so both must produce the same sort.
func TestSortItemsEmptyAuthorEqualsMissing(t *testing.T) {
	newer := time.Date(2026, time.January, 2, 0, 0, 0, 0, time.UTC)
	older := newer.AddDate(0, 0, -1)
	feed := func(author *gofeed.Person) []Item {
		return []Item{
			{Item: &gofeed.Item{Title: "Special", Author: author, PublishedParsed: &older}},
			{Item: &gofeed.Item{Title: "Missing", PublishedParsed: &newer}},
			{Item: &gofeed.Item{Title: "Bob", Author: &gofeed.Person{Name: "Bob"}}},
			{Item: &gofeed.Item{Title: "Alice", Author: &gofeed.Person{Name: "Alice"}}},
		}
	}

	withMissing := sortItems(feed(nil), "author")
	withEmpty := sortItems(feed(&gofeed.Person{Name: ""}), "author")

	for i := range withMissing {
		if withMissing[i].Title != withEmpty[i].Title {
			t.Errorf("item %d: missing author gives %q, empty author gives %q", i, withMissing[i].Title, withEmpty[i].Title)
		}
	}
}
