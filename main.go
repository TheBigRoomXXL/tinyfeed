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
