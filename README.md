canary
======

like ping, but for http

[`canaryio/canary`](https://github.com/canaryio/canary) is the spiritual successor to [`canaryio/sensord`](https://github.com/canaryio/sensord) and [`canaryio/canaryd`](https://github.com/canaryio/canaryd).

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
* duration of request / response in milliseconds
* was the response judged as healthy
* (optional) error message if the response was unhealthy
