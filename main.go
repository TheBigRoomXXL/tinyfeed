package main

import (
	"embed"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
)

//lint:ignore U1000 go:embed built-in.html
var htmlTemplate embed.FS

// flags
var limit int
var name string
var stylesheet string
var templatePath string

var rootCmd = &cobra.Command{
	Use:   "tinyfeed [FEED_URL ...]",
	Short: "Aggregate a collection of feed into static HTML page",
	Long:  "Aggregate a collection of feed into static HTML page. Only RSS, Atom and JSON feeds are supported.",
	Example: `  single feed      tinyfeed lovergne.dev/rss.xml > index.html
  multiple feeds   cat feeds.txt | tinyfeed > index.html`,
	Args: cobra.ArbitraryArgs,
	Run:  tinyfeed,
}

func init() {
	rootCmd.Flags().IntVarP(
		&limit,
		"limit",
		"l",
		50,
		"How many articles will be included",
	)
	rootCmd.Flags().StringVarP(
		&name,
		"name",
		"n",
		"Feed",
		"Name of the aggregated feed.",
	)
	rootCmd.Flags().StringVarP(
		&stylesheet,
		"stylesheet",
		"s",
		"",
		"Path to an external CSS stylesheet",
	)
	rootCmd.Flags().StringVarP(
		&templatePath,
		"template",
		"t",
		"built-in",
		"Path to a custom html+go template file.",
	)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func tinyfeed(cmd *cobra.Command, args []string) {
	strdinArgs, err := stdinToArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	args = append(args, strdinArgs...)

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "you must input at list one feed url.")
		return
	}

	feeds := []*gofeed.Feed{}
	items := []*gofeed.Item{}
	fp := gofeed.NewParser()
	for _, url := range args {
		feed, err := fp.ParseURL(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse feed: %s\n", err)
			continue
		}
		feeds = append(feeds, feed)
		items = append(items, feed.Items...)
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].PublishedParsed.After(*items[j].PublishedParsed)
	})

	items = items[0:min(len(items), limit)]

	err = printHTML(feeds, items)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func Preview(item *gofeed.Item) string {
	if len(item.Description) > 0 {
		return truncstr(item.Description, 600)
	} else {
		return truncstr(item.Content, 600)
	}
}

func Domain(item *gofeed.Item) string {
	url, err := url.Parse(item.Link)
	if err != nil {
		log.Fatal(err)
	}
	hostname := strings.TrimPrefix(url.Hostname(), "www.")
	return hostname
}

func printHTML(feeds []*gofeed.Feed, items []*gofeed.Item) error {
	ts, err := template.New(templatePath).
		Funcs(template.FuncMap{"Preview": Preview, "Domain": Domain}).
		ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("error loading html template: %s", err)
	}

	data := struct {
		Metadata map[string]string
		Items    []*gofeed.Item
		Feeds    []*gofeed.Feed
	}{
		Metadata: map[string]string{"name": name, "stylesheet": stylesheet},
		Items:    items,
		Feeds:    feeds,
	}

	err = ts.Execute(os.Stdout, data)
	if err != nil {
		return fmt.Errorf("error rendering html template: %s", err)
	}

	return nil
}
