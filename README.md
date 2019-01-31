# Magic: The Gathering Arena Log Parser #
[![Build Status](https://travis-ci.org/gathering-gg/parser.svg?branch=master)](https://travis-ci.org/gathering-gg/parser)
[![GoDoc](https://godoc.org/github.com/gathering-gg/parser?status.svg)](https://godoc.org/github.com/gathering-gg/parser)

_Parser_ is a Go library and command line interface to parse the
`output_log.txt` created by MTGA. Cross platform with minimal dependencies, the
executable is _small_. The CLI is currently used to send the parsed data to
[Gathering.gg](https://gathering.gg), and you need an account to use it. Future
work can be done to parse the data to local files.

## Usage ##
To use the parser, you'll need a [Go dev
environment](https://golang.org/doc/install). To use the parser locally:

```
$ go get github.com/gathering-gg/parser
```

If you want to install the CLI and use that:
```
$ go install github.com/gathering-gg/parser/cli
```
You can then execute it by running (You may need to add [Go's bin to your
path](https://github.com/golang/go/wiki/GOPATH)):
```
$ gathering -token=YOUR_GATHERING_GG_TOKEN
```
You can get your token from [Gathering.gg](https://gathering.gg).

#### Binaries ####
If you just want to execute a binary, you can get them from the
[releases](https://github.com/gathering-gg/parser/releases) page.

## Development ##
Once you have your Go dev environment setup, fetch the library:

```
$ go get github.com/gathering-gg/parser
$ cd $GOPATH/src/github.com/gathering-gg/parser
```

You can the start hacking on the library. To build the library, you run `go
build`, however there are some compile time variables that should be set. The
current command used to build is:

```
go build -ldflags "-X 'github.com/gathering-gg/gathering/config.Root=https://api.gathering.gg' -X 'github.com/gathering-gg/gathering/config.Version=0.0.1'" -o gathering ./cli
```
You can also build for specific platforms by appending `GOOS` and `GOARCH`
environment variables.

Tests:
```
$ go test -v ./...
```


## Changelog ##
See the [releases](https://github.com/gathering-gg/parser/releases) page.

## Contributing ##
Please contribute! If you have a problem with the parser you can [open an
issue](https://github.com/gathering-gg/parser/issues/new). Feel free to ask for
new features as well!

## Future Work ##
1. Performance. Not a lot of work has gone into performance - I'd like this to be the fastest parser around.

