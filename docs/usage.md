# Usage

!!! info
    This page use real feeds URL so that you can copy-paste the command and see the result. If you one the links is dead, please open an issue on the [repository](https://github.com/TheBigRoomXXL/tinyfeed/issues)

## Basic usage

Start by displaying the help message to check that tinyfeed is correctly installed
```bash
tinyfeed --help
```

Then this is how you use tinyfeed in it's most basic form:
```bash
tinyfeed --output "index.html" https://lovergne.dev/rss
```
This tell tinyfeed to fetch the feed at `https://lovergne.dev/rss`, generate a webpage for it and output the result to `index.html`. The value of `--output` can be a relative or absolute path. 

You can then open the resulting `index.html` in your browser using the `file` scheme with an url like `file:///absolute/page/to/your/index.html`. Or, on most OS, you can right-click on the `index.html` in your file explorer and use "Open With` and select your browser. 


If you want to process multiple feeds you can pass multiple URLs to tinyfeed:
```bash
tinyfeed --output "index.html" https://lovergne.dev/rss https://blog.codingconfessions.com/feed
```
Or use an input file containing all of your feeds (preferable if you have manny):
```bash
tinyfeed --input feeds.txt --output index.html
```
Example of `feeds.txt` file:
```txt
# Software Engineering
https://lovergne.dev/rss.xml

# Cyber Security
https://words.filippo.io/rss/
https://feeds.feedburner.com/TroyHunt

# Frontend
https://tonsky.me/atom.xml
https://andy-bell.co.uk/feed.xml
https://www.htmhell.dev/feed.xml
```

Notice that lines starting with `#` are comments and won't be parsed by tinyfeed.

## Configuration

You can change the default behavior of tinyfeed and customize it's settings using
flags:
```txt
Flags:
  -D, --daemon               Whether to execute the program in a daemon mode.
  -d, --description string   Add a description after the name of your page
  -h, --help                 help for tinyfeed
  -i, --input string         Path to a file with a list of feeds.
  -I, --interval int         Duration in minutes between execution. Ignored if not in daemon mode. (default 1440)
  -l, --limit int            How many articles to display (default 256)
  -n, --name string          Title of the page. (default "Feed")
  -o, --output string        Path to a file to save the output to.
  -q, --quiet                Add this flag to silence warnings.
  -r, --requests int         How many simulaneous requests can be made (default 16)
  -s, --stylesheet string    Path to an external CSS stylesheet
  -t, --template string      Path to a custom HTML+Go template file.
  -T, --timeout int          timeout to get feeds in seconds (default 15)
```

Most flags are used to change default values but more advanced flags like `--stylesheet`, `--template` and `--daemon` have there own dedicated documentation page:

 - LIST_DOCUMENTSTION





