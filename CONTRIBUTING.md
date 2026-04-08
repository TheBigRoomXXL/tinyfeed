# Contributing to tinyfeed

If you are reading this document, it's probably because you want to contribute to tinyfeed, thank you, it's very nice!

This document will provide you with 
- rules to respect so that we can collaborate without issue
- guidance for developing in this repo (dev setup, tests, build, etc.)

## Rules

### 1. All contributors are welcome

Everybody is welcome, especially beginners, tinyfeed is a very small codebase so it's a good
place to learn. 


### 2. Not all contributions are welcome

I am very opinionated about tinyfeed, I have my own philosophy for this project and I will
enforce it. Most notably tinyfeed is a minimalist project, that's why it's called tiny, 
because of that there are lots of features that are possible, and even easy to implement, 
but that I don't want. However, in tinyfeed I try to make exceptions possible for users that
are ready to tinker through customization (styling, scripting, templating), so if I don't 
want a feature in the core, I will probably give you some hints on how to circumvent that.

If you want to contribute but you don't know what to do, look at the [`contribution-welcome` issues](https://github.com/TheBigRoomXXL/tinyfeed/issues?q=is%3Aissue%20state%3Aopen%20label%3Acontribution-welcome)

### 3. Discussion before contribution

Before you contribute a feature I want you to open an issue, or respond to one, to make it 
known to me that you want to contribute and what you want to contribute. I will answer
with an approbation or a refusal and some specific guidelines about what you are trying to
achieve. This saves time for everyone, I don't want to tell you "I don't want it" after you
have done a lot of work to implement something. 

For fixes and documentation you can contribute directly but don't hesitate to open an issue
if you have questions or need guidance.


### 4. No breaking changes

I won't break existing user setups for any feature, so no breaking changes!

tinyfeed is not a library, but it still has an interface to respect in order to avoid breaking changes. It is composed of:
- The CLI: which commands and options can be called by the user
- The schema of the data exposed to templating


### 5. Be kind and have a nice time

Hey, it's an open source project, if we don't have a nice time, what's the point?


## Development

### Fork and Pull Request

For contributions we use the fork and pull request workflow. If you are not familiar with it, you can find a nice guide here:
- [GitHub Standard Fork & Pull Request Workflow](https://gist.github.com/Chaser324/ce0505fbed06b947d962 )

### Setting yourself up for development

To develop on tinyfeed all you need is the Go programming language installed on your computer (version > 1.23).

Download the dependencies:
```sh
go mod download
```

Run the latest version of the source code:
```sh
go run . --help # Equivalent to tinyfeed --help
```

Run the test suite:
```sh
go test -vet=all -v -race */*.go
```

Build the binary for release (this is done by the CI):
```sh
CGO_ENABLED=0 go build -ldflags "-s -w"
```

Build the Docker image:
```sh
docker build .
```

### Coding guidelines
- Keep it simple, stupid
- Always try to add a test
- When writing tests, favor table-driven tests to cover more cases
- Don't just test the happy path, test the failure cases.
