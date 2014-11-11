canary
======

like ping, but for http

## Installation

```sh
$ go get github.com/canaryio/canary/cmd/canary
```

## Usage

```sh
$ canary http://www.canary.io
2014/12/24 21:23:12 http://www.canary.io 200 129 true
2014/12/24 21:23:13 http://www.canary.io 200 91 true
2014/12/24 21:23:14 http://www.canary.io 200 89 true
2014/12/24 21:23:15 http://www.canary.io 200 88 true
2014/12/24 21:23:16 http://www.canary.io 200 94 true
^C
```

## Output

The following fields are emitted:

* date
* time
* url
* http status code
* milliseconds
* healthy?
* (optional) error message
