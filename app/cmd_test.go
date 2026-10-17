package tinyfeed

import (
	"errors"
	"flag"
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

func TestValidateOrderBy(t *testing.T) {
	tests := []struct {
		label   string
		input   string
		wantErr bool
	}{
		{
			label:   "Valid publication-date",
			input:   "publication-date",
			wantErr: false,
		},
		{
			label:   "Valid author",
			input:   "author",
			wantErr: false,
		},
		{
			label:   "Valid update date",
			input:   "update-date",
			wantErr: false,
		},
		{
			label:   "Valid feed-name",
			input:   "feed-name",
			wantErr: false,
		},
		{
			label:   "Empty string",
			input:   "",
			wantErr: true,
		},
		{
			label:   "Case sensitivity check",
			input:   "FEED-NAME",
			wantErr: true,
		},
		{
			label:   "Invalid value",
			input:   "description",
			wantErr: true,
		},
		{
			label:   "Partial match",
			input:   "publication",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.label, func(t *testing.T) {
			err := validateOrderBy(test.input)

			if (err != nil) != test.wantErr {
				t.Errorf("validateOrderBy() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestParseCmdHelp(t *testing.T) {
	for _, arg := range []string{"-h", "-help", "--help"} {
		t.Run(arg, func(t *testing.T) {
			_, err := parseCmd([]string{arg}, nil)
			if !errors.Is(err, flag.ErrHelp) {
				t.Errorf("parseCmd() error = %v, want flag.ErrHelp", err)
			}
		})
	}
}

func TestParseFlagsToTheEnd(t *testing.T) {
	tests := []struct {
		label    string
		args     []string
		wantArgs []string
		wantOut  string
	}{
		{"no args", []string{}, []string{}, ""},
		{"flag before", []string{"-o", "out", "a"}, []string{"a"}, "out"},
		{"flag after", []string{"a", "-o", "out"}, []string{"a"}, "out"},
		{"flag between", []string{"a", "-o", "out", "b"}, []string{"a", "b"}, "out"},
		{"double dash", []string{"a", "--", "-o", "out"}, []string{"a", "-o", "out"}, ""},
	}

	for _, test := range tests {
		t.Run(test.label, func(t *testing.T) {
			var out string
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			fs.StringVar(&out, "o", "", "")

			got, err := parseFlagsToTheEnd(fs, test.args)
			if err != nil {
				t.Fatalf("parseFlagsToTheEnd() error = %v", err)
			}
			if !slices.Equal(got, test.wantArgs) {
				t.Errorf("parseFlagsToTheEnd() args = %v, want %v", got, test.wantArgs)
			}
			if out != test.wantOut {
				t.Errorf("parseFlagsToTheEnd() -o = %q, want %q", out, test.wantOut)
			}
		})
	}
}

func TestParseCmdStdin(t *testing.T) {
	stdin := strings.NewReader("https://b.example/rss\n# comment\n")
	config, err := parseCmd([]string{"https://a.example/rss"}, stdin)
	if err != nil {
		t.Fatalf("parseCmd() error = %v", err)
	}
	want := []string{"https://a.example/rss", "https://b.example/rss"}
	if !slices.Equal(config.Urls, want) {
		t.Errorf("parseCmd() Urls = %v, want %v", config.Urls, want)
	}
}

func TestParseCmdValidatesOrderBy(t *testing.T) {
	_, err := parseCmd([]string{"--order-by", "invalid", "https://a.example/rss"}, nil)
	if err == nil {
		t.Error("parseCmd() error = nil, want error for invalid --order-by")
	}

	config, err := parseCmd([]string{"--order-by", "feed-name", "https://a.example/rss"}, nil)
	if err != nil {
		t.Fatalf("parseCmd() error = %v", err)
	}
	if config.OrderBy != "feed-name" {
		t.Errorf("parseCmd() OrderBy = %q, want %q", config.OrderBy, "feed-name")
	}
}
