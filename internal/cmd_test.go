package tinyfeed

import (
	"slices"
	"strings"
	"testing"
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
		{"http://foo.com http://foo.com", []string{"http://foo.com"}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			r := strings.NewReader(test.input)
			got, err := readerToArgs(r)
			if err != nil {
				t.Errorf("readerToArgs() error = %v", err)
				return
			}

			slices.Sort(test.want)
			slices.Sort(got)
			for i := range got {
				if got[i] != test.want[i] {
					t.Errorf("readerToArgs() got = %v, want %v", got, test.want)
					return
				}
			}
		})
	}
}
