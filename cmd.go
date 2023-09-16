package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "tinyfeed [FEED_URL ...]",
	Short: "Aggregate a collection of feed into static HTML page",
	Long:  "Aggregate a collection of feed into static HTML page. Only RSS, Atom and JSON feeds are supported.",
	Example: `  single feed      tinyfeed lovergne.dev/rss.xml > index.html
  multiple feeds   cat feeds.txt | tinyfeed > index.html`,
	Args:         cobra.ArbitraryArgs,
	RunE:         tinyfeed,
	SilenceUsage: true,
}

// flags
var limit int
var name string
var description string
var imageAllowed bool
var quiet bool
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
		"Title of the page.",
	)
	rootCmd.Flags().StringVarP(
		&description,
		"description",
		"d",
		"",
		"Add a description after the name of your page",
	)
	rootCmd.Flags().BoolVarP(
		&imageAllowed,
		"--images",
		"i",
		false,
		"Add this flag to load images in summaries",
	)
	rootCmd.Flags().BoolVarP(
		&quiet,
		"quiet",
		"q",
		false,
		"Add this flag to silence warnings.",
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
