
## Pipelining

!!! info
    This page is addressed to users with a good understanding of the terminal

tinyfeed is designed to work with basic pipelining. In its default behavior, tinyfeed expects a list of whitespace-separated feed URLs as arguments and will output HTML directly to `stdout`. It can take its input from `stdin` and use redirection for its output. tinyfeed also provides flags to use files directly. This provides a variety of options to manage I/O.

Given the following `feeds.txt` file:

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
Then the following commands are equivalent:
```bash
tinyfeed https://lovergne.dev/rss.xml \
    https://words.filippo.io/rss/ \
    https://feeds.feedburner.com/TroyHunt \
    https://tonsky.me/atom.xml \
    https://andy-bell.co.uk/feed.xml > index.html
```
```bash
cat feeds.txt | tinyfeed > index.html 
```
```bash
cat feeds.txt | tinyfeed -o index.html 
```
```bash
tinyfeed -i feeds.txt > index.html 
```
```bash
tinyfeed -i feeds.txt -o index.html
```

Which command you prefer will depend on your workflow but be aware that when using a redirection directly, like in the example, your HTML page will be blank while tinyfeed is processing and it will also stay blank if there is an error. To avoid that, use a temporary file or the `--output` flag: 

```bash
# This fix the blank page issue
cat feeds | tinyfeed > /tmp/tinyfeed && mv /tmp/tinyfeed /path/to/index.html
# OR
cat feeds | tinyfeed -o /path/to/index.html
```

!!! warning
    Commenting lines with `#` is only supported in input stream and input file. Adding `#` in the argument list directly will only cause a parsing error and will output a warning.
