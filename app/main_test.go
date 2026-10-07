package tinyfeed

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRunOrderByAuthorWithoutAuthors(t *testing.T) {
	feeds := []struct {
		name string
		body string
	}{
		{
			name: "RSS",
			body: `<rss version="2.0"><channel><title>Feed</title><link>https://example.com</link><description>Feed</description>
<item><title>Older</title><link>https://example.com/older</link><pubDate>Thu, 01 Jan 2026 00:00:00 +0000</pubDate></item>
<item><title>Newer</title><link>https://example.com/newer</link><pubDate>Fri, 02 Jan 2026 00:00:00 +0000</pubDate></item>
</channel></rss>`,
		},
		{
			name: "Atom",
			body: `<feed xmlns="http://www.w3.org/2005/Atom"><title>Feed</title><id>https://example.com</id><updated>2026-01-02T00:00:00Z</updated>
<entry><title>Older</title><id>https://example.com/older</id><published>2026-01-01T00:00:00Z</published><updated>2026-01-01T00:00:00Z</updated></entry>
<entry><title>Newer</title><id>https://example.com/newer</id><published>2026-01-02T00:00:00Z</published><updated>2026-01-02T00:00:00Z</updated></entry>
</feed>`,
		},
	}
	for _, feed := range feeds {
		t.Run(feed.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, feed.body)
			}))
			defer server.Close()
			dir := t.TempDir()
			templatePath := filepath.Join(dir, "items.gtpl")
			if err := os.WriteFile(templatePath, []byte("{{range .Items}}{{.Title}}\n{{end}}"), 0600); err != nil {
				t.Fatal(err)
			}
			outputPath := filepath.Join(dir, "index.html")
			config := &Config{
				Urls: []string{server.URL}, Limit: 10, LimitPerFeed: 10,
				Timeout: 5, RequestSemaphore: 1, OrderBy: "author",
				TemplatePath: templatePath, Output: outputPath,
			}
			if err := Run(config); err != nil {
				t.Fatal(err)
			}
			output, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatal(err)
			}
			if got, want := string(output), "Newer\nOlder\n"; got != want {
				t.Errorf("output = %q, want %q", got, want)
			}
		})
	}
}
