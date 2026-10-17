package tinyfeed

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tinyfeed "github.com/TheBigRoomXXL/tinyfeed/app"
)

var SNAPS_DIRECTORY = "snapshots/"
var FEEDS_DIRECTORY = "testsdata/"

var server *httptest.Server

func TestMain(m *testing.M) {
	server = httptest.NewServer(http.FileServer(http.Dir(FEEDS_DIRECTORY)))
	code := m.Run()
	server.Close()
	os.Exit(code)
}

// List of all feeds avaialble during tests
func Feeds() []string {
	entries, err := os.ReadDir(FEEDS_DIRECTORY)
	if err != nil {
		panic(err)
	}
	var feeds []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || strings.HasPrefix(name, ".") {
			continue
		}
		if ext := filepath.Ext(name); ext != ".xml" && ext != ".json" {
			continue
		}
		feeds = append(feeds, name)
	}
	return feeds
}

// Helper for snapshot based testing.
// Unlike traditional snapshot framework we always update
func Snapshot(t *testing.T, result []byte) {
	t.Helper()

	snap_path := SNAPS_DIRECTORY + t.Name() + ".html"
	defer func() {
		os.MkdirAll(filepath.Dir(snap_path), 0o755)
		os.WriteFile(snap_path, result, 0644)
	}()

	snapshot, err := os.ReadFile(snap_path)
	if err != nil {
		t.Fatal(err)
	}

	if string(snapshot) != string(result) {
		t.Fatalf("result does not match snapshot")
	}
}

func BaseConfig() tinyfeed.Config {
	c := tinyfeed.DefaultConfig
	c.Nonce = "very-stable-nounce"
	c.UpdatedAt = time.Date(2006, 1, 2, 3, 4, 5, 6, time.UTC)
	return c
}

func TestRunSingleFeed(t *testing.T) {
	feeds := Feeds()

	for _, feed := range feeds {
		t.Run(feed, func(t *testing.T) {
			dir := t.TempDir()
			outputPath := filepath.Join(dir, "index.html")

			config := BaseConfig()
			config.Urls = []string{server.URL + "/" + feed}
			config.Output = outputPath

			err := tinyfeed.Run(&config)

			if err != nil {
				t.Fatal(err)
			}

			output, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatal(err)
			}

			Snapshot(t, output)
		})
	}
}

func TestRunAllFeeds(t *testing.T) {
	feeds := Feeds()
	feedsUrl := make([]string, len(feeds))
	for i, feed := range feeds {
		feedsUrl[i] = server.URL + "/" + feed
	}

	sorts := []string{"publication-date", "update-date", "feed-name", "author"}
	for _, sort := range sorts {
		t.Run(sort, func(t *testing.T) {
			dir := t.TempDir()
			outputPath := filepath.Join(dir, "index.html")

			config := BaseConfig()
			config.Urls = feedsUrl
			config.Output = outputPath

			err := tinyfeed.Run(&config)

			if err != nil {
				t.Fatal(err)
			}

			output, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatal(err)
			}

			Snapshot(t, output)
		})
	}
}
