package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"
	"text/template"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
)

//go:embed built-in
var builtInTemplate string

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func tinyfeed(cmd *cobra.Command, args []string) error {
	strdinArgs, err := stdinToArgs()
	if err != nil {
		return fmt.Errorf("fail to parse stdin: %w", err)
	}

	args = append(args, strdinArgs...)

	if len(args) == 0 {
		return fmt.Errorf("no argument found, you must input at least one feed url. Use `tinyfeed --help` for examples")
	}

	feeds := parseFeeds(args)
	items := prepareItems(feeds)

	err = printHTML(feeds, items)
	if err != nil {
		return fmt.Errorf("fail to output HTML: %w", err)
	}
	return nil
}

func parseFeeds(url_list []string) []*gofeed.Feed {
	var sem = make(chan struct{}, requestSemaphore)
	var results = make(chan *gofeed.Feed)
	var wg sync.WaitGroup
	wg.Add(len(url_list))

	fp := gofeed.NewParser()
	fp.Client = &http.Client{Timeout: time.Duration(timeout * int(time.Second))}

	for _, url := range url_list {
		go func(url string) {
			defer func() {
				wg.Done()
				<-sem
			}()
			sem <- struct{}{}
			results <- parseFeed(url, fp)
		}(url)
	}

	go func() {
		wg.Wait()
		close(sem)
		close(results)
	}()

	feeds := []*gofeed.Feed{}
	for feed := range results {
		if feed != nil {
			feeds = append(feeds, feed)
		}
	}
	return feeds
}

func parseFeed(url string, fp *gofeed.Parser) *gofeed.Feed {
	feed, err := fp.ParseURL(url)
	if err != nil && !quiet {
		fmt.Fprintf(os.Stderr, "WARNING: fail to parse feed at %s: %s\n", url, err)
		return nil
	}

	return feed
}

func prepareItems(feeds []*gofeed.Feed) []*gofeed.Item {
	items := []*gofeed.Item{}
	for _, feed := range feeds {
		items = append(items, feed.Items...)
	}

	for i := 0; i < len(items); i++ {
		if items[i].Title == "" {
			items[i].Title = "<i>Untitled</i>"
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].PublishedParsed == nil {
			return false
		}
		if items[j].PublishedParsed == nil {
			return true
		}
		return items[i].PublishedParsed.After(*items[j].PublishedParsed)
	})

	return items[0:min(len(items), limit)]
}

func printHTML(feeds []*gofeed.Feed, items []*gofeed.Item) error {
	var err error
	var ts *template.Template

	if templatePath == "" {
		ts, err = template.New("built-in").
			Funcs(template.FuncMap{"domain": domain, "publication": publication}).
			Parse(builtInTemplate)
	} else {
		ts, err = template.New(templatePath).
			Funcs(template.FuncMap{"domain": domain, "publication": publication}).
			ParseFiles(templatePath)
	}
	if err != nil {
		return fmt.Errorf("fail to load HTML template: %w", err)
	}

	data := struct {
		Metadata map[string]string
		Items    []*gofeed.Item
		Feeds    []*gofeed.Feed
	}{
		Metadata: map[string]string{
			"name":        name,
			"description": description,
			"stylesheet":  stylesheet,
			"nonce":       randStr(20),
		},
		Items: items,
		Feeds: feeds,
	}

	err = ts.Execute(os.Stdout, data)
	if err != nil {
		return fmt.Errorf("fail to render HTML template: %w", err)
	}

	return nil
}
