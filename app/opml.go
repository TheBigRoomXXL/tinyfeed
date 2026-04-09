package tinyfeed

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type opmlDoc struct {
	XMLName  xml.Name      `xml:"opml"`
	Outlines []opmlOutline `xml:"body>outline"`
}

type opmlOutline struct {
	XMLUrl   string        `xml:"xmlUrl,attr"`
	Outlines []opmlOutline `xml:"outline"`
}

func collectFeedURLs(outlines []opmlOutline, out *[]string) {
	for _, o := range outlines {
		if o.XMLUrl != "" {
			*out = append(*out, o.XMLUrl)
		}
		collectFeedURLs(o.Outlines, out)
	}
}

func parseOPML(r io.Reader) ([]string, error) {
	var doc opmlDoc
	if err := xml.NewDecoder(r).Decode(&doc); err != nil {
		return nil, fmt.Errorf("invalid OPML: %w", err)
	}
	var urls []string
	collectFeedURLs(doc.Outlines, &urls)
	return urls, nil
}

func isOPMLContent(data []byte) bool {
	s := strings.ToLower(strings.TrimSpace(string(data)))
	idx := strings.Index(s, "<opml")
	return idx >= 0 && idx < 1024
}

// expandOPMLFile checks if a local file path is an OPML file and returns its feed URLs.
// Returns nil if the file is not OPML (so caller falls back to line-by-line parsing).
func expandOPMLFile(path string) ([]string, error) {
	if path == "" {
		return nil, nil
	}
	if !strings.HasSuffix(strings.ToLower(path), ".opml") {
		return nil, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if !isOPMLContent(data) {
		return nil, nil
	}
	return parseOPML(strings.NewReader(string(data)))
}

// expandOPMLURL fetches a URL, and if it's an OPML document returns the feed URLs inside.
// Returns nil, nil if not OPML (so caller keeps the original URL for gofeed).
func expandOPMLURL(rawURL string, client *http.Client) ([]string, error) {
	resp, err := client.Get(rawURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ct := strings.ToLower(resp.Header.Get("Content-Type"))
	couldBeOPML := ct == "" ||
		strings.Contains(ct, "opml") ||
		strings.Contains(ct, "text/xml") ||
		strings.Contains(ct, "application/xml")

	if couldBeOPML && isOPMLContent(data) {
		return parseOPML(strings.NewReader(string(data)))
	}
	return nil, nil
}

// resolveURLs expands any OPML URLs in the list to their contained feed URLs.
func resolveURLs(urls []string, timeout int) []string {
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	out := make([]string, 0, len(urls))

	for _, u := range urls {
		lower := strings.ToLower(u)
		if !strings.HasPrefix(lower, "http://") && !strings.HasPrefix(lower, "https://") {
			out = append(out, u)
			continue
		}

		// only bother fetching if content-type might be OPML (check extension or keyword)
		if strings.HasSuffix(lower, ".opml") || strings.Contains(lower, "opml") {
			expanded, err := expandOPMLURL(u, client)
			if err == nil && expanded != nil {
				out = append(out, expanded...)
				continue
			}
		}
		out = append(out, u)
	}
	return out
}
