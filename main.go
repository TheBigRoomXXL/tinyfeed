package main

import (
	_ "embed"
	"fmt"
	"html"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
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
	log.SetOutput(os.Stderr)

	// We get the inputs stdin at the start to that it can be reused by the daemon
	strdinArgs, err := stdinToArgs()
	if err != nil {
		return fmt.Errorf("fail to parse stdin: %w", err)
	}
	args = append(args, strdinArgs...)

	err = run(args)

	if !daemon || err != nil {
		return err
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	ticker := time.NewTicker(time.Minute * time.Duration(interval))

	for {
		select {
		case <-signalChan:
			return nil
		case <-ticker.C:
			err = run(args)
			if err != nil {
				return err
			}
		}
	}
}

func run(args []string) error {
	// We append inputs from file here so that it can be updated without
	// reloading the daemon
	inputArgs, err := fileToArgs(input)
	if err != nil {
		return fmt.Errorf("fail to parse input file: %w", err)
	}
	args = append(args, inputArgs...)

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
	if url[0] =='#'{
		log.Printf("Ignoring Commented URL: %s\n", url[1:])
		return nil
	}
	feed, err := fp.ParseURL(url)
	if err != nil && !quiet {
		log.Printf("WARNING: fail to parse feed at %s: %s\n", url, err)
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
			items[i].Title = "Untitled"
		}

		// Some stirng are already html escaped ans when they are parsed by
		// html/template it create a double escape so we must unescape first.
		items[i].Title = html.UnescapeString(items[i].Title)
		items[i].Link = html.UnescapeString(items[i].Link)
		items[i].Published = html.UnescapeString(items[i].Published)
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
	currDate := time.Now()
	data := struct {
		Metadata map[string]string
		Items    []*gofeed.Item
		Feeds    []*gofeed.Feed
	}{
		Metadata: map[string]string{
			"name":        name,
			"description": description,
			"stylesheet":  stylesheet,
			"nonce":       generateNonce(256),
			"day":currDate.Weekday().String(),
			"datetime":currDate.Format(time.DateTime),
		},
		Items: items,
		Feeds: feeds,
	}

	var outFile io.WriteCloser
	if output != "" {
		outFile, err = os.Create(output)
		if err != nil {
			return err
		}
		defer outFile.Close()
	} else {
		outFile = os.Stdout
	}

	err = ts.Execute(outFile, data)
	if err != nil {
		return fmt.Errorf("fail to render HTML template: %w", err)
	}

	return nil
}
