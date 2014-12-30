canary
======

a small cli tool to help you measure availability and performance of a given website.

It's like ping, but for HTTP.

## Installation

```sh
$ go get github.com/canaryio/canary/cmd/canary
```

## Usage

```sh
$ canary http://www.canary.io
2014-12-28T14:44:32Z http://www.canary.io 200 96 true
2014-12-28T14:44:33Z http://www.canary.io 200 92 true
2014-12-28T14:44:34Z http://www.canary.io 200 89 true
2014-12-28T14:44:35Z http://www.canary.io 200 124 true
^C
```

## Output

The following fields are emitted:

* date
* time
* url
* http status code
* duration of request / response in milliseconds
* was the response judged as healthy
* (optional) error message if the response was unhealthy
