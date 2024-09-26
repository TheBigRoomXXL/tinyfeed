package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "tinyfeed [FEED_URL ...]",
	Short: "Aggregate a collection of feed into static HTML page",
	Example: `  single feed      tinyfeed lovergne.dev/rss.xml > index.html
  multiple feeds   cat feeds.txt | tinyfeed > index.html
  daemon mode      tinyfeed --daemon -i feeds.txt -o index.html `,

	Args:         cobra.ArbitraryArgs,
	RunE:         tinyfeed,
	SilenceUsage: true,
}

// flags
var limit int
var timeout int
var requestSemaphore int
var name string
var description string
var quiet bool
var stylesheet string
var templatePath string
var input string
var output string
var daemon bool
var interval int64

func init() {
	rootCmd.Flags().IntVarP(
		&limit,
		"limit",
		"l",
		256,
		"How many articles to display",
	)
	rootCmd.Flags().IntVarP(
		&requestSemaphore,
		"requests",
		"r",
		16,
		"How many simulaneous requests can be made",
	)
	rootCmd.Flags().IntVarP(
		&timeout,
		"timeout",
		"T",
		15,
		"timeout to get feeds in seconds",
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
	rootCmd.Flags().StringVarP(
		&input,
		"input",
		"i",
		"",
		"Path to a file with a list of feeds.",
	)
	rootCmd.Flags().StringVarP(
		&output,
		"output",
		"o",
		"",
		"Path to a file to save the output to.",
	)
	rootCmd.Flags().BoolVarP(
		&daemon,
		"daemon",
		"D",
		false,
		"Whether to execute the program in a daemon mode.",
	)
	rootCmd.Flags().Int64VarP(
		&interval,
		"interval",
		"I",
		1440,
		"Duration in minutes between execution. Ignored if not in daemon mode.",
	)
}
