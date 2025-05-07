package tinyfeed

import (
	"testing"

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
			item := &gofeed.Item{Link: test.input}
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
			item := &gofeed.Item{Link: test}
			got := domain(item)
			if got != "" {
				t.Errorf("domain() got = %v, want empty string", got)
			}
		})
	}
}
