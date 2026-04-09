package tinyfeed

import (
	"os"
	"strings"
	"testing"
)

const subscriptionListOPML = `<?xml version="1.0" encoding="UTF-8"?>
<opml version="2.0">
  <head>
    <title>Subscriptions</title>
  </head>
  <body>
    <outline text="News" title="News">
      <outline type="rss" text="NY Times" xmlUrl="http://www.nytimes.com/services/xml/rss/nyt/Technology.xml" />
      <outline type="rss" text="Scripting News" xmlUrl="http://www.scripting.com/rss.xml" />
    </outline>
    <outline type="rss" text="Wired" xmlUrl="http://www.wired.com/news_drop/netcenter/netcenter.rdf" />
  </body>
</opml>`

const emptyOPML = `<?xml version="1.0" encoding="UTF-8"?>
<opml version="2.0">
  <head><title>Empty</title></head>
  <body></body>
</opml>`

const noXmlUrlOPML = `<?xml version="1.0" encoding="UTF-8"?>
<opml version="2.0">
  <head><title>Categories only</title></head>
  <body>
    <outline text="News" title="News">
      <outline text="No feed here" />
    </outline>
  </body>
</opml>`

func TestParseOPML(t *testing.T) {
	tests := []struct {
		label   string
		input   string
		want    []string
		wantErr bool
	}{
		{
			label: "nested and flat outlines",
			input: subscriptionListOPML,
			want: []string{
				"http://www.nytimes.com/services/xml/rss/nyt/Technology.xml",
				"http://www.scripting.com/rss.xml",
				"http://www.wired.com/news_drop/netcenter/netcenter.rdf",
			},
		},
		{
			label: "empty body",
			input: emptyOPML,
			want:  []string{},
		},
		{
			label: "outlines without xmlUrl are skipped",
			input: noXmlUrlOPML,
			want:  []string{},
		},
		{
			label:   "invalid XML",
			input:   "<not valid xml",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			got, err := parseOPML(strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseOPML() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if len(got) != len(tt.want) {
				t.Fatalf("parseOPML() len = %d, want %d; got %v", len(got), len(tt.want), got)
			}
			for i, u := range got {
				if u != tt.want[i] {
					t.Errorf("parseOPML()[%d] = %q, want %q", i, u, tt.want[i])
				}
			}
		})
	}
}

func TestExpandOPMLFile(t *testing.T) {
	t.Run("valid .opml file", func(t *testing.T) {
		f, err := os.CreateTemp("", "test-*.opml")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		f.WriteString(subscriptionListOPML)
		f.Close()

		urls, err := expandOPMLFile(f.Name())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(urls) != 3 {
			t.Errorf("expected 3 urls, got %d: %v", len(urls), urls)
		}
	})

	t.Run("non-opml extension returns nil", func(t *testing.T) {
		f, err := os.CreateTemp("", "feeds-*.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		f.WriteString("http://example.com/feed.xml\n")
		f.Close()

		urls, err := expandOPMLFile(f.Name())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if urls != nil {
			t.Errorf("expected nil for non-opml file, got %v", urls)
		}
	})

	t.Run("empty path returns nil", func(t *testing.T) {
		urls, err := expandOPMLFile("")
		if err != nil || urls != nil {
			t.Errorf("expected nil, nil; got %v, %v", urls, err)
		}
	})
}

func TestIsOPMLContent(t *testing.T) {
	tests := []struct {
		label string
		input string
		want  bool
	}{
		{"opml tag present", `<?xml version="1.0"?><opml version="2.0">`, true},
		{"opml tag with whitespace", "  \n<opml>", true},
		{"rss feed", `<?xml version="1.0"?><rss version="2.0">`, false},
		{"atom feed", `<?xml version="1.0"?><feed xmlns="http://www.w3.org/2005/Atom">`, false},
		{"empty", "", false},
		{"opml tag buried deep", strings.Repeat("x", 1025) + "<opml>", false},
	}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			got := isOPMLContent([]byte(tt.input))
			if got != tt.want {
				t.Errorf("isOPMLContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
