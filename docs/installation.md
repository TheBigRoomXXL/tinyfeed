
# Installation

## Install from binary 

You can download the official binaries from the [releases page](https://github.com/TheBigRoomXXL/tinyfeed/releases/latest/). tinyfeed is built completely statically so there are no dependencies to install. Currently arm64 and amd64 architecture on Linux, Mac and Windows and FreeBSD is supported. 

If you need something an architecture or OS currently not support, please open an issue and I will add it to the releases process if it's supported by golang cross-compilation.

Here is a one-liner to install the binary for Linux amd64:

```bash
sudo wget https://github.com/TheBigRoomXXL/tinyfeed/releases/latest/download/tinyfeed_linux_arm64 -O /usr/local/bin/tinyfeed
sudo chmod +x /usr/local/bin/tinyfeed
tinyfeed --help
```
To adapt this script, you can just change `_linux_arm64` to your system.

!!! info
    Binary for Linux and Windows and compressed using the awesome [upx](https://upx.github.io/) which reduce there size by around 50%!

## Install with Go

```bash 
go install github.com/TheBigRoomXXL/tinyfeed@latest
```

## Install with docker

```bash
docker run thebigroomxxl/tinyfeed --help
```

