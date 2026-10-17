package tinyfeed

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"html"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"sync"
	"syscall"
	"time"

	"github.com/mmcdole/gofeed"
)

//go:embed templates/index.html.gtpl
var builtInTemplate string

//go:embed templates/index.opml.gtpl
var opmlTemplate string

type Item struct {
	*gofeed.Item
	FeedName string
}

func Main() {
	config, err := parseCmd(os.Args[1:], pipedStdin())
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			printHelp()
			os.Exit(0)
		}
		log.Printf("%s\n", err)
		os.Exit(1)
	}

	err = Run(config)
	if err != nil {
		log.Printf("%s\n", err)
		os.Exit(1)
	}

	if !config.Daemon {
		return
	}

	// Daemon mode
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	ticker := time.NewTicker(time.Minute * time.Duration(config.Interval))

	for {
		select {
		case <-signalChan:
			return
		case <-ticker.C:
			err = Run(config)
			if err != nil {
				log.Printf("%s\n", err)
				os.Exit(1)
			}
		}
	}
}

func Run(config *Config) error {
	// We append inputs from file here so that it can be updated without reloading the daemon
	fileUrls, err := fileToArgs(config.Input)
	if err != nil {
		return fmt.Errorf("fail to parse input file: %w", err)
	}
	urls := append(config.Urls, fileUrls...)

	if len(urls) == 0 {
		return fmt.Errorf(
			"no argument found, you must input at least one feed url. Use `tinyfeed --help` for examples",
		)
	}

	feeds := parseFeeds(urls, config)
	items := prepareItems(feeds, config)

	err = printHTML(feeds, items, config)
	if err != nil {
		return fmt.Errorf("fail to output HTML: %w", err)
	}
	return nil
}

func parseFeeds(url_list []string, config *Config) []*gofeed.Feed {
	sem := make(chan struct{}, config.RequestSemaphore)
	results := make([]*gofeed.Feed, len(url_list))
	var wg sync.WaitGroup

	fp := gofeed.NewParser()
	fp.UserAgent = "tinyfeed/v1"
	fp.Client = &http.Client{Timeout: time.Duration(config.Timeout * int(time.Second))}

	for i, url := range url_list {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			results[i] = parseFeed(url, fp, config)
		}()
	}
	wg.Wait()

	// Discard failed fetch
	success := slices.DeleteFunc(results, func(f *gofeed.Feed) bool { return f == nil })

	return success
}

func parseFeed(url string, fp *gofeed.Parser, config *Config) *gofeed.Feed {
	feed, err := fp.ParseURL(url)
	if err != nil {
		if !config.Quiet {
			log.Printf("WARNING: fail to process feed at %s: %s\n", url, err)
		}
		return nil
	}

	if feed.FeedLink == "" {
		feed.FeedLink = url
	}

	sortFeedItems(feed, config.OrderBy)
	feed.Items = feed.Items[:min(len(feed.Items), config.LimitPerFeed)]

	return feed
}

func prepareItems(feeds []*gofeed.Feed, config *Config) []Item {
	var items []Item

	for _, feed := range feeds {
		for _, item := range feed.Items {
			items = append(items, Item{
				item,
				feed.Title,
			})
		}
	}

	for i := range items {
		if items[i].Title == "" {
			items[i].Title = "Untitled"
		}

		// Some string are already html escaped inside the feeds and when
		// html/template Run it re-escape them, creating double escape. In
		// order to avoid malformed string we must unescape first.
		items[i].Title = html.UnescapeString(items[i].Title)
		items[i].Link = html.UnescapeString(items[i].Link)
		items[i].Description = html.UnescapeString(items[i].Description)
		items[i].Content = html.UnescapeString(items[i].Content)
		items[i].Published = html.UnescapeString(items[i].Published)
		items[i].FeedName = html.UnescapeString(items[i].FeedName)
	}

	items = sortItems(items, config.OrderBy)

	return items[0:min(len(items), config.Limit)]
}

func printHTML(feeds []*gofeed.Feed, items []Item, config *Config) error {
	var err error
	var ts *template.Template
	var templateFuncs = template.FuncMap{"domain": domain, "publication": publication, "opmlRSSVersion": opmlRSSVersion}
	switch config.TemplatePath {
	case "":
		ts, err = template.New("templates/index.html.gtpl").
			Funcs(templateFuncs).
			Parse(builtInTemplate)
	case "opml":
		ts, err = template.New("templates/index.opml.gtpl").
			Funcs(templateFuncs).
			Parse(opmlTemplate)
	default:
		ts, err = template.New(filepath.Base(config.TemplatePath)).
			Funcs(templateFuncs).
			ParseFiles(config.TemplatePath)
	}

	if err != nil {
		return fmt.Errorf("fail to load HTML template: %w", err)
	}
	currDate := config.UpdatedAt
	if currDate.IsZero() {
		currDate = time.Now()
	}

	nonce := config.Nonce
	if nonce == "" {
		nonce = generateNonce(256)
	}

	data := struct {
		// Metadata is an arbitrary map, there is no backward compatibility promise on it's content.
		// Changing it's content won't be concidered breaking change. (But we keep it as stable as possible)
		// In the future, if we are confident that it won't evolve, we will replace it with a struct
		// to ensure backward compatibiliy
		Metadata    map[string]string
		Items       []Item
		Feeds       []*gofeed.Feed
		Stylesheets []string
		Scripts     []string
	}{
		Metadata: map[string]string{
			"name":            config.Name,
			"description":     config.Description,
			"nonce":           nonce,
			"day":             currDate.Weekday().String(),
			"datetime":        currDate.Format(time.DateTime),
			"datetimeRFC3339": currDate.Format(time.RFC3339),
			"datetimeRFC822":  currDate.Format(time.RFC822),
		},
		Items:       items,
		Feeds:       feeds,
		Stylesheets: config.Stylesheets,
		Scripts:     config.Scripts,
	}

	var outFile io.WriteCloser
	if config.Output != "" {
		outFile, err = os.Create(config.Output)
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
