package main

import (
	"strings"
	"testing"

	"github.com/mmcdole/gofeed"
)

func TestReaderToArgs(t *testing.T) {
	var tests = []struct {
		input string
		want  []string
	}{
		{"", []string{}},
		{"# comment\n", []string{}},
		{"# comment\nhttp://foo.com\n", []string{"http://foo.com"}},
		{"# comment\nhttp://foo.com\nhttp://bar.net\n", []string{"http://foo.com", "http://bar.net"}},
		{"# comment\nhttp://foo.com\nhttp://bar.net\n# comment\n", []string{"http://foo.com", "http://bar.net"}},
		{"# comment\nhttp://foo.com\nhttp://bar.net\n# comment\nhttp://baz.org\n", []string{"http://foo.com", "http://bar.net", "http://baz.org"}},
		{"#comment\nhttp://foo.com\n", []string{"http://foo.com"}},
		{" #comment\nhttp://foo.com\n", []string{"http://foo.com"}},
		{" # comment\nhttp://foo.com\n", []string{"http://foo.com"}},
		{"\n\nhttp://foo.com\n", []string{"http://foo.com"}},
		{"\n   \nhttp://foo.com\n", []string{"http://foo.com"}},
		{"\n \nhttp://foo.com\n", []string{"http://foo.com"}},
		{"http://foo.com", []string{"http://foo.com"}},
		{"http://foo.com http://baz.org", []string{"http://foo.com", "http://baz.org"}},
		{" http://foo.com ", []string{"http://foo.com"}},
		{"# http://foo.com ", []string{""}},
		{" # http://foo.com ", []string{""}},
		{"\n # http://foo.com ", []string{""}},
		{"http://foo.com\thttp://baz.org", []string{"http://foo.com", "http://baz.org"}},
		{"http://foo.com http://baz.org", []string{"http://foo.com", "http://baz.org"}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			r := strings.NewReader(test.input)
			got, err := readerToArgs(r)
			if err != nil {
				t.Errorf("readerToArgs() error = %v", err)
				return
			}
			for i := range got {
				if got[i] != test.want[i] {
					t.Errorf("readerToArgs() got = %v, want %v", got, test.want)
					return
				}
			}
		})
	}
}

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
