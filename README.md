canary
======

like ping, but for http

## Installation

```sh
$ go get github.com/canaryio/canary
```

## Usage

```sh
# let's measure http://www.canary.io
$ canary http://www.canary.io
2014-11-27T18:43:16Z resh.local http://www.canary.io 23.235.40.133 0 200 439.199000 481.048000 547.217000 596.852000
2014-11-27T18:43:17Z resh.local http://www.canary.io 23.235.40.133 0 200 0.060000 0.060000 44.395000 89.276000
2014-11-27T18:43:18Z resh.local http://www.canary.io 23.235.40.133 0 200 0.046000 0.047000 49.809000 93.389000
2014-11-27T18:43:19Z resh.local http://www.canary.io 23.235.40.133 0 200 0.023000 0.023000 42.613000 87.549000
^C
```

## Output

The following fields are emitted:

* timestamp
* hostname of machine running the commmand
* url being measured
* ip address connected to
* libcurl status code
* http status code
* dns lookup time
* time to connect
* time to first byte
* total time

