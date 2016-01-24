go-retry
========

A retry command by golang on CLI.

[![Travis](https://img.shields.io/travis/linyows/go-retry.svg?style=flat-square)][travis]
[![GitHub release](http://img.shields.io/github/release/linyows/go-retry.svg?style=flat-square)][release]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[travis]: https://travis-ci.org/linyows/go-retry
[release]: https://github.com/linyows/go-retry/releases
[license]: https://github.com/linyows/go-retry/blob/master/LICENSE
[godocs]: http://godoc.org/github.com/linyows/go-retry

Description
-----------

Retry n times with interval for your command until to zero for exit-status.

Usage
-----

```sh
$ retry -i 5s -c 2 /usr/lib64/nagios/plugins/check_http -w 10 -c 15 -H localhost
```

Install
-------

To install, use `go get`:

```bash
$ go get -d github.com/linyows/go-retry
```

Contribution
------------

1. Fork ([https://github.com/linyows/go-retry/fork](https://github.com/linyows/go-retry/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

Author
------

[linyows](https://github.com/linyows)
