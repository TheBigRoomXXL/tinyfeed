# tinyfeed 

![banner](.images/banner.svg)

**tinyfeed** is a CLI tool that generate a static HTML page from a collection of feeds.

It's dead simple, no database, no config file, just a CLI and some HTML 

Give it a list of RSS, Atom or JSON feeds urls and it will generate a single HTML page for
it. Then you can effortlessly set it up in `crond` `systemd` or `openrc` and voilà, you’ve
got yourself an webpage that aggregate your favorite feeds.

## Feature

- RSS, Atom and JSON feeds are all supported thanks to the awesome 
[gofeed library](https://github.com/mmcdole/gofeed)
- Highly customizable, especially with the ability to use external stylesheet and templates.
- Dark / Light theme based on system preference
- Generated page is lightweight and fully accessible
- Support a daemon mode to re-generate the output periodically


## Demo and Screenshots

### Live demo: [feed.lovergne.dev](https://feed.lovergne.dev/)

Visited links are in yellow, unvisited in blue. 
![screenshots of feed.lovergne.dev](/.images/screenshots.png)


## Usage

The CLI app is design to work with basic pipelining and stdout redirections.

tinyfeed expect a list of whitespace separated feeds urls as argument. Lines can 
be commented out with a `#` at the start. 

```
Usage:
  tinyfeed [FEED_URL ...] [flags]

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
  -s, --stylesheet string    Path to an external CSS stylesheet
  -t, --template string      Path to a custom HTML+Go template file.
  -T, --timeout int          Timeout to get feeds in seconds (default 15)
```

⚠️ When using a redirection directly, like in the example, your HTML page will be
blank while tinyfeed is processing and it will also stay blank if there is an error.
 To avoid that, use a tempory file or the --output flag: 

```bash
cat feeds | tinyfeed > /tmp/tinyfeed && mv /tmp/tinyfeed /path/to/index.html
# OR
cat feeds | tinyfeed -o /path/to/index.html
```

Example of input file:
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


## Installation

### Install from binary 

You can download the official binaries from the [releases page](https://github.com/TheBigRoomXXL/tinyfeed/releases/latest/). Currently arm64 and amd64 architecture on Linux, Mac and Windows and FreeBSD is supported. If you need something else than that, please open an issue and I will add it to the releases process if it's supported by golang 
cross-compilation.

Here is a quick example of how to install the binary for linux:

```bash
wget https://github.com/TheBigRoomXXL/tinyfeed/releases/latest/download/tinyfeed_linux_arm64
chmod +x tinyfeed_linux_arm64
sudo mv tinyfeed_linux_arm64 /usr/local/bin/tinyfeed
tinyfeed --help
```

### Install with Go

```bash 
go install github.com/TheBigRoomXXL/tinyfeed@latest
```

### Install with docker

```bash
docker run thebigroomxxl/tinyfeed --help
```


## Recipes

### Docker

This is a simple example of how to run **tinyfeed** in a docker container. This will mount
an entire directory, if you want to bind only the input/output files instead you will have
use [bind mounts](https://docs.docker.com/engine/storage/bind-mounts/).

```bash
docker run --restart unless-stopped  -v /your/path:/app thebigroomxxl/tinyfeed --daemon -i feeds.txt -o index.html
```

Docker compose equivalent:

```yaml
services:
    tinyfeed:
        image: thebigroomxxl/tinyfeed
        command: --daemon -i feeds.txt -o index.html
        volumes:
            - /path/to/your/feeds/:/app
        restart: unless-stopped
```


### Systemd

Here is a simple systemd service file to run **tinyfeed** when your system start and 
update the page every 12 hours. With this setup you can edit the feeds list at 
`~/feeds.txt` and the output will be updated after the next run. 

With the output file at `~/index.html` you can access it locally at `file:///home/<USER>/index.html`.
If, instead, you want to serve it with a web server you can save it in the web server root directory.

```ini
# /etc/systemd/system/tinyfeed.service

[Unit]
Description=tinyfeed service
After=network.target

[Service]
Type=simple
Restart=always
User=<USER>
WorkingDirectory=/home/<USER>/
ExecStart=/usr/local/bin/tinyfeed --daemon -i feeds.txt -o index.html -I 720

[Install]
WantedBy=mutli-user.target
```

If you have SELinux enabled you will need to allow systemd to execute binaries in the `usr/local/bin` directory with the following commands:
```bash
sudo semanage fcontext -a -t bin_t /usr/local/bin 
sudo chcon -Rv -u system_u -t bin_t /usr/local/bin 
sudo restorecon -R -v /usr/local/bin
```

### OpenRC

To create an OpenRC service, you just need to add an init file at `/etc/init.d/tinyfeed`.
Then, you can enable the service with `rc-update add tinyfeed default` and start or stop it with `rc-service tinyfeed <COMMAND>`. Below, you can find the most minimal service file that will enable you to start and stop tinyfeed with the feed list and rendered page located at `/etc/tinyfeed`.

```ini
#!/sbin/openrc-run

depend() {
	need net
	use dns 
}

command="/usr/local/bin/tinyfeed"
command_args="--daemon -i /etc/tinyfeed/feeds.txt -o /etc/tinyfeed/index.html"
command_background=true
pidfile="/run/${RC_SVCNAME}.pid"
```
For more advanced patterns (like running as your user instead of root), you can check out the official [OpenRC documentation](https://github.com/OpenRC/openrc/blob/master/service-script-guide.md) on the subject.

### Github Action + Github Page

You can use GitHub Actions to periodically generate an updated page, similar to a cron job, and host it using GitHub Pages. The advantage of this approach is that it’s free and serverless. Additionally, you can directly update your feed list from GitHub.

To use this method you will need to create a github repository with:
- GitHub Pages enabled
- A [repository variable](https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/store-information-in-variables) named FEEDS that contains the list of feeds you want to aggregate
- The following GitHub Actions workflow file located at `.github/workflows/daily.yml`:

```yaml
name: update-demo-daily
on:
  schedule:
    - cron: "0 6 * * *" # every day at 6am
  push:
    tags:
      - "refresh-github-page" # to manually trigger a refresh

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: prepare-page
        env:
          FEEDS: ${{vars.FEEDS}}
        run: |
          wget -q -O tinyfeed https://github.com/TheBigRoomXXL/tinyfeed/releases/latest/download/tinyfeed_linux_arm64
          chmod +x tinyfeed
          mkdir www
          echo $FEEDS | ./tinyfeed > ./www/index.html
          cp www/index.html www/404.html 
        # 404.html allows every path to be served by the index.html

      - name: upload-page
        uses: actions/upload-pages-artifact # add @ to specify a version
        with:
          path: www/

  deploy:
    needs: build
    permissions:
      pages: write # to deploy to Pages
      id-token: write # to verify the deployment originates from an appropriate source
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    steps:
      - name: deploy-to-github-pages
        id: deployment
        uses: actions/deploy-pages # add @ to specify a version
```

## External HTML+Go template 

You can provide you own template for page generation. For an exemple template
check out the [built-in one](https://github.com/TheBigRoomXXL/tinyfeed/blob/main/built-in).
To learn about HTML+Go template check the [official documentation](https://pkg.go.dev/html/template). 

Inside you template you will have access to data with the following struct and functions:

```go
type data struct {
    Metadata map[string]string
    Items    []*gofeed.Item
    Feeds    []*gofeed.Feed
}

func publication(item *gofeed.Item) string

func domain(item *gofeed.Item) string
```


## Feedback, help or bug report

You have created a page with tinyfeed and you want to share it? You can open a
merge request or an issue to add it to the demo section. 

If you need anything related to this project wether it's' just giving feedback,
help to understand something or feature request just open a issue on this repos.


## Acknowledgement

The project was heavily inspired by the awesomely simple [tinystatus](https://github.com/bderenzo/tinystatus)
and message boards like Lobste.rs and Hacker News.

Thank you @MariaLetta for the awesome [free-gophers-pack ](https://github.com/MariaLetta/free-gophers-pack)
wich I adapted for the banner.
