package main

import (
	"embed"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
)

//go:embed template.html
var htmlTemplate embed.FS

var rootCmd = &cobra.Command{
	Use:   "tinyfeed [FEED_URL ...]",
	Short: "Aggregate a collection of feed into static HTML page",
	Long:  "Aggregate a collection of feed into static HTML page. Only RSS, Atom and JSON feeds are supported.",
	Example: `  single feed      tinyfeed lovergne.dev/rss.xml > index.html
  multiple feeds   cat feeds.txt | tinyfeed > index.html`,
	Args: cobra.ArbitraryArgs,
	Run:  tinyfeed,
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
		fmt.Fprintln(os.Stderr, "You must input at list one feed url.")
		return
	}

	items := []*gofeed.Item{}
	fp := gofeed.NewParser()
	for _, url := range args {
		feed, _ := fp.ParseURL(url)
		items = append(items, feed.Items...)
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].PublishedParsed.After(*items[j].PublishedParsed)
	})

	items = items[0:min(len(items), 49)]

	err = printHTML(items)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func printHTML(items []*gofeed.Item) error {
	ts, err := template.ParseFiles("template.html")
	if err != nil {
		return fmt.Errorf("error loading html template: %s", err)
	}

	err = ts.Execute(os.Stdout, items)
	if err != nil {
		return fmt.Errorf("error rendering html template: %s", err)
	}

	return nil
}

func stdinToArgs() ([]string, error) {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("error parsing stdin: %s", err)
	}

	return strings.Fields(string(input)), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
