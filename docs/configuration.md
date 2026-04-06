# Configuration

## Flags

tinyfeed is configured using command line flags. Flags should be input before the feed list.
Each flag has a long and a short form.

You can use the help flag `--help` to display the full list of available flags, their description and their default: 
flags:
```txt
❯ tinyfeed --help
Aggregate a collection of feeds into static HTML page

Usage:
  tinyfeed [flags] [FEED_URL ...]

Examples:
  single feed      tinyfeed lovergne.dev/rss.xml > index.html
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

For the full tinyfeed manual, please visit: https://feed.lovergne.dev/
```

This page of the manual will only focus on the configuration flags.


## Examples

### Change the time between refresh in daemon mode

Refresh every hour:
```bash
# Long form
tinyfeed --daemon --input feeds.txt --output index.html --interval 60

# Short form
tinyfeed -D -i feeds.txt -o index.html -I 60
```

Refresh every week:
```bash
# Long form
tinyfeed --daemon --input feeds.txt --output index.html --interval 10080

# Short form
tinyfeed -D -i feeds.txt -o index.html -I 10080
```

### Paginate the results

Limit results to only 20 items in total:
```bash
# Long form
tinyfeed --input feeds.txt --output index.html --limit 20

# Short form
tinyfeed -i feeds.txt -o index.html -l 20
```

Limit results to only 5 items per feed:
```bash
# Long form
tinyfeed --input feeds.txt --output index.html --limit-per-feed 5

# Short form
tinyfeed -i feeds.txt -o index.html -L 5
```

Limit results to only 5 items per feed and 100 items in total:
```bash
# Long form
tinyfeed --input feeds.txt --output index.html --limit 100 --limit-per-feed 5

# Short form
tinyfeed -i feeds.txt -o index.html -l 100 -L 5
```

Limit to 5 items per feed and sort by feed name:
```bash
# Long form
tinyfeed --input feeds.txt --output index.html --limit-per-feed 5 --order-by feed-name

# Short form
tinyfeed -i feeds.txt -o index.html -L 5 -O feed-name
```

### Disable warnings

This option can be useful in automated scripts or cronjobs to only have errors in `stderr`.
```bash
# Long form
tinyfeed --input feeds.txt --output index.html --quiet

# Short form
tinyfeed -i feeds.txt -o index.html -q
```

### Change how feeds are fetched

Decrease the number of concurrent requests to 4:
```bash
# Long form
tinyfeed --input feeds.txt --output index.html --requests 4

# Short form
tinyfeed -i feeds.txt -o index.html -r 4
```

Change the timeout duration on individual requests to 60 seconds:
```bash
# Long form
tinyfeed --input feeds.txt --output index.html --timeout 60

# Short form
tinyfeed -i feeds.txt -o index.html -T 60
```
