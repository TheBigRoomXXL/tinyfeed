# Docker

!!! Info
    Docker is a tool that is used to automate the deployment of applications in lightweight containers so that applications can work efficiently in different environments in isolation.  
    - *from [Wikipedia](https://en.wikipedia.org/wiki/Docker_(software))*

The official Docker image for tinyfeed is [thebigroomxxl/tinyfeed](https://hub.docker.com/r/thebigroomxxl/tinyfeed)

## Running with pipelining

!!! warning
    You will not be able to use [pipelining](pipelining.md) if you don't use Docker's `-i` flag (interactive mode).

In this mode, you pass the feeds as arguments to the command and redirect the output using pipelining. This method is simpler because it does not require mounting a volume.

Single feed:
```bash
docker run -i thebigroomxxl/tinyfeed https://lovergne.dev/rss  > index.html
```

Multiple feeds from a file :
```bash
cat feeds.txt | docker run -i thebigroomxxl/tinyfeed > index.html
```

## Running with input/output files as a volume


!!! warning 
    Files for the `--input` and `--output` flags won't be available if you don't mount them.

This mounts an entire directory where your input (`feeds.txt`) and output (`index.html`) files are located. If you want to bind only the input/output files instead, you will need to use [bind mounts](https://docs.docker.com/engine/storage/bind-mounts/).

```bash
docker run -v /your/path:/app thebigroomxxl/tinyfeed -i feeds.txt -o index.html
```

Docker Compose equivalent:

```yaml
services:
  tinyfeed:
    image: thebigroomxxl/tinyfeed
    command: -i feeds.txt -o index.html
    volumes:
      - /path/to/your/feeds/:/app
```

## Running as a service (daemon mode)

tinyfeed has a [daemon mode](daemon.md) where it continuously runs and periodically updates the output file. To use it, pass the `--daemon` flag and set it to restart automatically using Docker's `--restart unless-stopped`.

This is the previous example adapted to run as a service:

```bash
docker run --restart unless-stopped -v /your/path:/app thebigroomxxl/tinyfeed --daemon -i feeds.txt -o index.html
```

Docker Compose equivalent:

```yaml
services:
  tinyfeed:
    image: thebigroomxxl/tinyfeed
    command: --daemon -i feeds.txt -o index.html
    volumes:
      - /path/to/your/feeds/:/app
    restart: unless-stopped
```

As always, tinyfeed only produces an HTML page—it does not run a web server. It's up to you to serve the generated HTML file with a secondary service like [Caddy](https://caddyserver.com/) or [NGINX](https://nginx.org/). Here is an example with Caddy (I prefer it because it is simpler to configure and handles HTTPS for you automatically):

```yaml
services:
  tinyfeed:
    image: thebigroomxxl/tinyfeed
    command: --daemon -i feeds.txt -o index.html
    volumes:
      - /path/to/your/feeds/:/app
    restart: unless-stopped

  http-server:
    image: caddy
    ports:
      - "80:80"
      - "443:443"
      - "443:443/udp"
    volumes:
      - /path/to/your/caddyfile:/etc/caddy/Caddyfile
      - /path/to/your/feeds/:/srv/feed
    restart: unless-stopped
```
With you caddy file being someting like:
```caddyfile
feeds.example.com {
	rewrite * index.html
	root * /srv/feed
	file_server
}
```
