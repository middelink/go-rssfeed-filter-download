# go-rssfeed-filter-download

[![GoDoc](https://godoc.org/github.com/middelink/go-rssfeed-filter-download?status.svg)](https://godoc.org/github.com/middelink/go-rssfeed-filter-download)
[![License](https://img.shields.io/github/license/middelink/go-rssfeed-filter-download.svg)](https://github.com/middelink/go-rssfeed-filter-download/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/middelink/go-rssfeed-filter-download.svg?branch=master)](https://travis-ci.org/middelink/go-rssfeed-filter-download)
[![Coverage Status](https://coveralls.io/repos/github/middelink/go-rssfeed-filter-download/badge.svg?branch=master)](https://coveralls.io/github/middelink/go-rssfeed-filter-download?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/middelink/go-rssfeed-filter-download)](https://goreportcard.com/report/github.com/middelink/go-rssfeed-filter-download)

> Extract torrent information from rss feed to sync up local volumes

Fetch a rss feed and extract a list of series. Then compare these series
against the QNAP Download area to see in which series we are interested.
For any new episodes found in the rss feed, start a download through QNAPs
Download Station.

## Why?

I am watching too many anime series to track them manually. Then when including
holidays and other times I'm interrupt, left me with series occasionally
missing an episode. Annoying! Hence this program.

## Usage

```shell
$ /usr/local/bin/go-rssfeed-filter-download --help
Usage of /usr/local/bin/go-rssfeed-filter-download:
  -baseuri string
    	url for the qnap, e.g. http://192.168.1.5:8080/
  -filter string
    	Which resolutions we do not care for (default "1080p,720p")
  -pass string
    	qnap pass to log in with, defaults to env QNAP_PASS
  -rss string
    	rss feed to monitor (default "http://horriblesubs.info/rss.php?res=sd")
  -silent
    	silence mode
  -user string
    	qnap user to log in as, defaults to env QNAP_USER
```

In fact I run mine from cron:
```shell
MAILTO=<your email>
QNAP_USER=<your user>
QNAP_PASS=<your password>

05 */8 * * * /usr/local/bin/go-rssfeed-filter-download --baseuri http://<your qnap>:8080/ --silent
```

## Install

### Automatic

```sh
$ go get github.com/middelink/go-rssfeed-filter-download
```

### Manual

First clone the repository.

```sh
$ git clone https://github.com/middelink/go-rssfeed-filter-download && cd go-rssfeed-filter-download
```

And run the command for installing the package.

```sh
$ go install .
```

## Contributing

Take a look at the open
[issues](https://github.com/middleink/go-rssfeed-filter-download/issues) and submit a PR!

## License

MIT © [Pauline Middelink](http://www.polyware.nl/~middelink)
