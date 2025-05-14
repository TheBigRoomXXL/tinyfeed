# Usage

!!! info
    This page uses real feeds URLs so that you can copy-paste the command and see the result. If one of the links is dead, please [open an issue].(https://github.com/TheBigRoomXXL/tinyfeed/issues)

## Basic usage

Start by displaying the help message to check that tinyfeed is correctly installed
```bash
tinyfeed --help
```

Then this is how you use tinyfeed in its most basic form:
```bash
tinyfeed --output "index.html" https://lovergne.dev/rss
```
This tell tinyfeed to fetch the feed at `https://lovergne.dev/rss`, generate a webpage for it and output the result to `index.html`. The value of `--output` can be a relative or absolute path. 

You can then open the resulting `index.html` in your browser using the `file` scheme with an URL like `file:///absolute/page/to/your/index.html`. Or, on most OS, you can right-click on the `index.html` in your file explorer and use "Open With" and select your browser. 

If you want to process multiple feeds, you can pass multiple URLs to tinyfeed:
```bash
tinyfeed --output "index.html" https://lovergne.dev/rss https://blog.codingconfessions.com/feed
```
Or use an input file containing all of your feeds (preferable if you have many):
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

## Docker caveat

If you are running tinyfeed through Docker, you should be aware of the following details:

- Files for the `--input` and `--output` flags won't be available if you don't mount them.
- You will not be able to use [pipelining](pipelining.md) if you don't use Docker's `-i` flag (interactive mode).

Single feed without mount:
```bash
docker run -i thebigroomxxl/tinyfeed https://lovergne.dev/rss > index.html
```

Multiple feeds from a file without mount:
```bash
cat feeds.txt | docker run -i thebigroomxxl/tinyfeed > index.html
```
For more guidance, check out the [Docker section](docker.md) of the documentation.

## Configuration

You can change the default behavior of tinyfeed and customize its settings using
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
  -s, --stylesheet string    Link to an external CSS stylesheet
  -t, --template string      Path to a custom HTML+Go template file.
  -T, --timeout int          timeout to get feeds in seconds (default 15)
```

Most flags are used to change default values, but more advanced flags like `--daemon`, `--stylesheet` and `--template` have their own dedicated documentation page:

- [Understanding daemon mode](daemon.md)
- [Custom theme with styling](styling.md/)
- [Custom layout with templating](templating.md)


## Hosting the tinyfeed webpage

tinyfeed is only a static site generator (even if it's a one-page website) so once you are happy with the page you generated you will need a way to host it and a way to update it.

For the hosting, the solution is an HTTP server like [NGINX](https://nginx.org/) or [Caddy](https://caddyserver.com/). If you are unfamiliar with HTTP servers, I would recommend [Caddy](https://caddyserver.com/) because it is simpler to configure and handles HTTPS for you automatically.

For the updating, you will have a lot of options and this will depend on your setup. For guidance on that subject, the workflow section of the documentation presents 5 options:

- [Cron](cron.md): running tinyfeed periodically in a cron job. This is the simplest setup.
- [Docker](docker.md): running tinyfeed as a containerized service. This is a good option if you are already familiar with Docker.
- [systemd](systemd.md) and [OpenRC](openrc.md): running tinyfeed as a service managed by your [init system](https://en.wikipedia.org/wiki/Init). This is the traditional (and dependencies-free) way to integrate a daemon with your OS.
- [GitHub Action and Page](github.md): run tinyfeed the serverless and free (as in free beer) way. Use GitHub Action to periodically generate your page like a cron job and GitHub Page to replace a traditional HTTP server.
