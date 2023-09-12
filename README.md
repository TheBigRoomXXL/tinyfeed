# tinyfeed 

**tinyfeed** is a CLI that generate a static HTML page from a collection of feeds.

You can effortlessly set it up in crond or systemd and direct your HTTP server 
to the generated `index.html` and voilà, you’ve got yourself a webpage 
that aggregate your feeds.

The design of the webpage is heavily inspired by hackernews.

## Demo

### [feed.lovergne.dev/](https://feed.lovergne.dev/)

## Feature

- RSS, Atom and JSON feeds are all supported thanks to the awesome 
[gofeed library](https://github.com/mmcdole/gofeed)
- Options to use external stylesheet and templates to configure it to your taste.
- Generated page is lightweight and fully accessible
- no exernal dependency


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

## Installations

### Install with Go:
```bash 
go install github.com/TheBigRoomXXL/tinyfeed@latest
```

### Install from binary 

If you need something else than `arm64` please open an issue and I will add it
to the release process.
```bash
wget https://github.com/TheBigRoomXXL/tinyfeed/releases/latest/download/tinyfeed_linux_arm64
chmod +x tinyfeed_linux_arm64
sudo mv tinyfeed_linux_arm64 /usr/local/bin/tinyfeed
tinyfeed --help
```

If you are on Alpine you also need to install [gcompat](https://pkgs.alpinelinux.org/package/edge/main/x86_64/gcompat) to fix the usual musl / glibc compatibility
issues.
```bash
apk add gcompat
```

## External HTML+Go template 

You can provide you own template for page generation. For an exemple template
check out the [built-in one](https://github.com/TheBigRoomXXL/tinyfeed/blob/main/.built-in).
To learn about HTML+Go template check the [official documentation](https://pkg.go.dev/html/template). 

Inside you template you will have access to data with the following struct:

```go
type data struct {
    Metadata map[string]string
    Items    []*gofeed.Item
    Feeds    []*gofeed.Feed
}
```
