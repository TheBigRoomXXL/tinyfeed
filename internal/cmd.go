package tinyfeed

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// implement flag.Value
// Adapted from https://stackoverflow.com/a/28323276
type StringRepeatable []string

func (i *StringRepeatable) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *StringRepeatable) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// flags
type Config struct {
	Urls             []string
	Limit            int
	LimitPerFeed     int
	Timeout          int
	RequestSemaphore int
	Name             string
	Description      string
	Quiet            bool
	Stylesheets      StringRepeatable
	Scripts          StringRepeatable
	TemplatePath     string
	Input            string
	Output           string
	Daemon           bool
	Interval         int64
	OrderBy          string
}

func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(0)
}

func printHelp() {
	log.Println(`Aggregate a collection of feeds into static HTML page

Usage:
  tinyfeed [flags] [FEED_URL ...]

Examples:
  single feed      tinyfeed https://feed.lovergne.dev/releases.atom > index.html
  multiple feeds   cat feeds.txt | tinyfeed > index.html
  daemon mode      tinyfeed --daemon -i feeds.txt -o index.html

Flags:

  Main flags:
  -i, --input string         Path to a file with a list of feeds.
  -o, --output string        Path to a file to save the output to.
  -D, --daemon               Whether to execute the program in a daemon mode.
  
  Customization flags:
  -n, --name string          Title of the page. (default "Feed")
  -d, --description string   Add a description after the name of your page
  -s, --stylesheet string    Link to an external CSS stylesheet
  -S, --script string        Link to an external JavaScript file
  -t, --template string      Path to a custom HTML+Go template file.

  Configuration flags:
  -I, --interval int         Duration in minutes between execution. Ignored if not in daemon mode. (default 1440)
  -l, --limit int            How many articles to display in total (default 256)
  -L, --limit-per-feed int   Maximum number of articles to display per feed (default 256)
  -q, --quiet                Add this flag to silence warnings.
  -r, --requests int         How many simultaneous requests can be made (default 16)
  -T, --timeout int          Timeout to get feeds in seconds (default 15)
  -O, --order-by string      How to order the articles. Accept 'publication-date', 'update-date', 'feed-name','author'. (default publication-date)

  -h, --help                 help for tinyfeed

For the full tinyfeed manual, please visit: https://feed.lovergne.dev/`)
	os.Exit(0)
}

func parseCmd() (*Config, error) {
	var config Config
	config.Scripts = make(StringRepeatable, 0)
	config.Stylesheets = make(StringRepeatable, 0)

	fs := flag.NewFlagSet("tinyfeed", flag.ContinueOnError)
	fs.Usage = func() {} // Disable default usage message

	fs.IntVar(&config.Limit, "limit", 256, "How many articles to display in total")
	fs.IntVar(&config.Limit, "l", 256, "How many articles to display in total")

	fs.IntVar(&config.LimitPerFeed, "limit-per-feed", 256, "Maximum number of articles to display per feed")
	fs.IntVar(&config.LimitPerFeed, "L", 256, "Maximum number of articles to display per feed")

	fs.IntVar(&config.RequestSemaphore, "requests", 16, "How many simulaneous requests can be made")
	fs.IntVar(&config.RequestSemaphore, "r", 16, "How many simulaneous requests can be made")

	fs.IntVar(&config.Timeout, "timeout", 15, "Timeout to get feeds in seconds")
	fs.IntVar(&config.Timeout, "T", 15, "Timeout to get feeds in seconds")

	fs.StringVar(&config.Name, "name", "Feed", "Title of the page")
	fs.StringVar(&config.Name, "n", "Feed", "Title of the page")

	fs.StringVar(&config.Description, "description", "", "Add a description after the name of your page")
	fs.StringVar(&config.Description, "d", "", "Add a description after the name of your page")

	fs.BoolVar(&config.Quiet, "quiet", false, "Silence warnings")
	fs.BoolVar(&config.Quiet, "q", false, "Silence warnings")

	fs.Var(&config.Stylesheets, "stylesheet", "Link to an external CSS stylesheet")
	fs.Var(&config.Stylesheets, "s", "Link to an external CSS stylesheet")

	fs.Var(&config.Scripts, "script", "Link to an external JavaScript file")
	fs.Var(&config.Scripts, "S", "Link to an external JavaScript file")

	fs.StringVar(&config.TemplatePath, "template", "", "Path to a custom HTML+Go template file")
	fs.StringVar(&config.TemplatePath, "t", "", "Path to a custom HTML+Go template file")

	fs.StringVar(&config.Input, "input", "", "Path to a file with a list of feeds")
	fs.StringVar(&config.Input, "i", "", "Path to a file with a list of feeds")

	fs.StringVar(&config.Output, "output", "", "Path to a file to save the output to")
	fs.StringVar(&config.Output, "o", "", "Path to a file to save the output to")

	fs.BoolVar(&config.Daemon, "daemon", false, "Whether to execute the program in a daemon mode")
	fs.BoolVar(&config.Daemon, "D", false, "Whether to execute the program in a daemon mode")

	fs.Int64Var(&config.Interval, "interval", 1440, "Duration in minutes between execution. Ignored if not in daemon mode")
	fs.Int64Var(&config.Interval, "I", 1440, "Duration in minutes between execution. Ignored if not in daemon mode")

	fs.StringVar(&config.OrderBy, "order-by", "publication-date", "How to order the articles")
	fs.StringVar(&config.OrderBy, "O", "publication-date", "How to order the articles")

	args, err := parseFlagsToTheEnd(fs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	// We get the inputs stdin at the start to that it can be reused by the daemon
	strdinArgs, err := stdinToArgs()
	if err != nil {
		return nil, fmt.Errorf("failed to parse stdin: %w", err)
	}
	config.Urls = append(args, strdinArgs...)

	err = validateOrderBy(config.OrderBy)
	if err != nil {
		return nil, err
	}

	return &config, nil
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

func validateOrderBy(orderBy string) error {
	accepted := []string{"publication-date", "update-date", "feed-name", "author"}
	for _, accept := range accepted {
		if orderBy == accept {
			return nil
		}
	}
	return fmt.Errorf(
		"Invalid value '%s' for flag --order-by, accepted values are %s.",
		orderBy,
		strings.Join(accepted, ", "),
	)
}
