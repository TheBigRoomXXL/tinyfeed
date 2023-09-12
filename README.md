# tinyfeed 

**tinyfeed** is a CLI that generate static HTML page from a collection of feeds.
Currently RSS, Atom and JSON feeds are all supported thanks to the awesome 
[gofeed library](https://github.com/mmcdole/gofeed)

You can effortlessly set it up in cron or systemd and direct your HTTP server 
to the generated `index.html` and voilà, you’ve got yourself a feed webpage.

## Usage

The CLI app is design to work with basic pipelining and stdout redirections. 

```
Usage:
  tinyfeed [flags] FEED_URL1 FEED_URL2 ...

Examples:
  single feed      tinyfeed lovergne.dev/rss.xml > index.html
  multiple feeds   cat feeds.txt | tinyfeed > index.html

Flags:
  -h, --help                help for tinyfeed
  -l, --limit int           How many articles will be included (default 50)
  -n, --name string         Name of the aggregated feed. (default "Feed")
  -s, --stylesheet string   Path to an external CSS stylesheet
  -t, --template string     Path to a custom HTML+Go template file.
```
