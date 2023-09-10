package main

import (
	"fmt"
	"io"
	"os"
	"strings"

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
	args = append(args, stdinToArgs()...)
	for _, arg := range args {
		fmt.Println(arg)
	}
}

func stdinToArgs() []string {
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	return strings.Fields(string(stdin))
}
