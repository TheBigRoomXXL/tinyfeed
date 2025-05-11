package tinyfeed

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// flags
var fs *flag.FlagSet
var limit int
var limitPerFeed int
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
	log.SetOutput(os.Stderr)
	log.SetFlags(0)

	fs = flag.NewFlagSet("tinyfeed", flag.ContinueOnError)
	fs.Usage = func() {} // Disable default usage message

	fs.IntVar(&limit, "limit", 256, "How many articles to display in total")
	fs.IntVar(&limit, "l", 256, "How many articles to display in total")

	fs.IntVar(&limitPerFeed, "limit-per-feed", 256, "Maximum number of articles to display per feed")
	fs.IntVar(&limitPerFeed, "L", 256, "Maximum number of articles to display per feed")

	fs.IntVar(&requestSemaphore, "requests", 16, "How many simulaneous requests can be made")
	fs.IntVar(&requestSemaphore, "r", 16, "How many simulaneous requests can be made")

	fs.IntVar(&timeout, "timeout", 15, "Timeout to get feeds in seconds")
	fs.IntVar(&timeout, "T", 15, "Timeout to get feeds in seconds")

	fs.StringVar(&name, "name", "Feed", "Title of the page")
	fs.StringVar(&name, "n", "Feed", "Title of the page")

	fs.StringVar(&description, "description", "", "Add a description after the name of your page")
	fs.StringVar(&description, "d", "", "Add a description after the name of your page")

	fs.BoolVar(&quiet, "quiet", false, "Silence warnings")
	fs.BoolVar(&quiet, "q", false, "Silence warnings")

	fs.StringVar(&stylesheet, "stylesheet", "", "Path to an external CSS stylesheet")
	fs.StringVar(&stylesheet, "s", "", "Path to an external CSS stylesheet")

	fs.StringVar(&templatePath, "template", "", "Path to a custom HTML+Go template file")
	fs.StringVar(&templatePath, "t", "", "Path to a custom HTML+Go template file")

	fs.StringVar(&input, "input", "", "Path to a file with a list of feeds")
	fs.StringVar(&input, "i", "", "Path to a file with a list of feeds")

	fs.StringVar(&output, "output", "", "Path to a file to save the output to")
	fs.StringVar(&output, "o", "", "Path to a file to save the output to")

	fs.BoolVar(&daemon, "daemon", false, "Whether to execute the program in a daemon mode")
	fs.BoolVar(&daemon, "D", false, "Whether to execute the program in a daemon mode")

	fs.Int64Var(&interval, "interval", 1440, "Duration in minutes between execution. Ignored if not in daemon mode")
	fs.Int64Var(&interval, "I", 1440, "Duration in minutes between execution. Ignored if not in daemon mode")

}

func printHelp() {
	log.Println(`Aggregate a collection of feed into static HTML page

Usage:
  tinyfeed [flags] [FEED_URL ...]

Examples:
  single feed      tinyfeed lovergne.dev/rss.xml > index.html
  multiple feeds   cat feeds.txt | tinyfeed > index.html
  daemon mode      tinyfeed --daemon -i feeds.txt -o index.html

Flags:
  -D, --daemon               Whether to execute the program in a daemon mode.
  -d, --description string   Add a description after the name of your page
  -h, --help                 help for tinyfeed
  -i, --input string         Path to a file with a list of feeds.
  -I, --interval int         Duration in minutes between execution. Ignored if not in daemon mode. (default 1440)
  -l, --limit int            How many articles to display in total (default 256)
  -L, --limit-per-feed int   Maximum number of articles to display per feed (default 256)
  -n, --name string          Title of the page. (default "Feed")
  -o, --output string        Path to a file to save the output to.
  -q, --quiet                Add this flag to silence warnings.
  -r, --requests int         How many simulaneous requests can be made (default 16)
  -s, --stylesheet string    Link to an external CSS stylesheet
  -t, --template string      Path to a custom HTML+Go template file.
  -T, --timeout int          Timeout to get feeds in seconds (default 15)

For more instructions on how to integrate tinyfeed with your workflow, please visit:
https://github.com/TheBigRoomXXL/tinyfeed#recipes`)
	os.Exit(0)
}

// This function allow to parse the flags that are passed after aguments in order to
// avoid breaking after the replacement of cobra by the std flag package.
// Basicly the following flag is parsed correclty using cobra
//
//	tinyfeed https://lovergne.dev/rss.xml -o index.html
//
// But not using the std flag package.
// Adapted from https://github.com/golang/go/issues/36744 by thinkerbot
func parseFlagsToTheEnd(fs *flag.FlagSet) ([]string, error) {
	err := fs.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}
	args := make([]string, 0)
	for i := len(os.Args) - len(fs.Args()) + 1; i < len(os.Args); {
		if i > 1 && os.Args[i-2] == "--" {
			break
		}
		args = append(args, fs.Arg(0))
		err = fs.Parse(os.Args[i:])
		if err != nil {
			return nil, err
		}
		i += 1 + len(os.Args[i:]) - len(fs.Args())
	}
	return append(args, fs.Args()...), nil
}

func stdinToArgs() ([]string, error) {
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return readerToArgs(os.Stdin)
	}
	return []string{}, nil
}

func fileToArgs(filepath string) ([]string, error) {
	if filepath == "" {
		return []string{}, nil
	}
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening input file: %s", err)
	}
	return readerToArgs(file)
}

func readerToArgs(reader io.Reader) ([]string, error) {
	input, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading input: %s", err)
	}
	lines := strings.Split(string(input), "\n")

	argsSet := make(map[string]struct{})
	for _, s := range lines {
		lineTrimmed := strings.TrimSpace(s)
		if strings.HasPrefix(lineTrimmed, "#") || lineTrimmed == "" {
			continue
		}
		for _, field := range strings.Fields(lineTrimmed) {
			argsSet[field] = struct{}{}
		}
	}

	i := 0
	args := make([]string, len(argsSet))
	for k := range argsSet {
		args[i] = k
		i++
	}
	return args, nil
}
