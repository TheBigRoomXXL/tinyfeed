package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tinyfeed [FEED_URL ...]",
	Short: "Generate a static HTML page from a collection of feeds.",
	Long:  "Generate a static HTML page from a collection of feeds. Only RSS, Atom and JSON feeds are supported.",
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
		fmt.Fprintln(os.Stderr, "Error parsing stdin: ", err)
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

	printHTML(items)
}

func printHTML(items []*gofeed.Item) {
	html := `<!DOCTYPE html>
	<html lang="en" dir="ltr" itemscope itemtype="https://schema.org/WebPage" prefix="og:http://ogp.me/ns#">
	
	<head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
	
		<title>Feed</title>
	
		<link rel="icon" type="image/svg+xml" href="/icon.svg" />
	
		<!-- General -->
		<meta name="application-name" content="tinyfeed" />
		<meta name="author" content="Sebastien LOVERGNE" />
		<meta name="description" content="tiny feed reader" />
		<meta name="referrer" content="strict-origin" />
	</head>
	<body>
	<h1>Feed</h1>`
	for _, item := range items {
		html += fmt.Sprintf("<a href=\"%s\"><h2>%s</h2><a>", item.Link, item.Title)
	}
	html += "</body></html>"
	fmt.Println(html)
}

func stdinToArgs() ([]string, error) {
	//Check if stdin is Used
	stdin := os.Stdin
	// fi, err := stdin.Stat()
	// if err != nil {
	// 	return []string{}, nil
	// }
	// size := fi.Size()
	// if size == 0 {
	// 	return []string{}, nil
	// }

	input, err := io.ReadAll(stdin)
	if err != nil {
		return nil, err
	}

	return strings.Fields(string(input)), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
