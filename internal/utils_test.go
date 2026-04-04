package tinyfeed

import (
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
)

func TestDomain(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"http://foo.com", "foo.com"},
		{"https://www.foo.com", "foo.com"},
		{"http://www.foo.com", "foo.com"},
		{"http://foo.com/bar/baz", "foo.com"},
		{"http://www.foo.com/bar/baz", "foo.com"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			// we need to initialize the embedded field explicitely
			item := Item{
				Item: &gofeed.Item{Link: test.input},
			}
			got := domain(item)
			if got != test.want {
				t.Errorf("domain() got = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDomainMalformed(t *testing.T) {
	var tests = []string{
		"foocom",
		"",
		"http://",
		" no a url",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			item := Item{
				Item: &gofeed.Item{Link: test},
			}
			got := domain(item)
			if got != "" {
				t.Errorf("domain() got = %v, want empty string", got)
			}
		})
	}
}

func TestPublicationDate(t *testing.T) {
	date1 := time.Date(2026, 4, 4, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		label           string
		published       string
		publishedParsed *time.Time
		want            string
	}{
		{
			label:           "Nominal: Use parsed date",
			published:       "Fri, 04 Apr 2026 08:46:00 +0000",
			publishedParsed: &date1,
			want:            "2026-04-04",
		},
		{
			label:           "Parsed is nil, use raw string",
			published:       "April 2026",
			publishedParsed: nil,
			want:            "April 2026",
		},
		{
			label:           "Both empty/nil",
			published:       "",
			publishedParsed: nil,
			want:            "Once upon a time",
		},
		{
			label:           "Parsed is nil, string is whitespace",
			published:       " ",
			publishedParsed: nil,
			want:            "Once upon a time",
		},
		{
			label:           "Parsed is nil, string is line return",
			published:       "\n",
			publishedParsed: nil,
			want:            "Once upon a time",
		},
		{
			label:           "Leap year date",
			published:       "2024-02-29",
			publishedParsed: &date2,
			want:            "2024-02-29",
		},
	}

	for _, test := range tests {
		t.Run(test.label, func(t *testing.T) {
			item := Item{
				Item: &gofeed.Item{
					Published:       test.published,
					PublishedParsed: test.publishedParsed,
				},
			}
			got := publication(item)
			if got != test.want {
				t.Errorf("%s: want %s, got %s", test.label, test.want, got)
			}
		})
	}
}

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
			orderBy = test.orderBy

			got := sortItems(test.input)

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

	orderBy = "publication-date"
	items := []Item{itemNil, itemValid}
	sorted := sortItems(items)

	if sorted[0].Title != "ValidDate" {
		t.Errorf("Expected ValidDate to come before NilDate in descending sort")
	}
}
