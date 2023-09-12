package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "tinyfeed [FEED_URL ...]",
	Short: "Aggregate a collection of feed into static HTML page",
	Long:  "Aggregate a collection of feed into static HTML page. Only RSS, Atom and JSON feeds are supported.",
	Example: `  single feed      tinyfeed lovergne.dev/rss.xml > index.html
  multiple feeds   cat feeds.txt | tinyfeed > index.html`,
	Args: cobra.ArbitraryArgs,
	Run:  tinyfeed,
}

// flags
var limit int
var name string
var stylesheet string
var templatePath string

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
		"",
		"Path to a custom HTML+Go template file.",
	)
}
