canary
======

a small cli tool to help you measure availability and performance of a given website.

It's like ping, but for HTTP.

## Installation

```sh
$ go get github.com/canaryio/canary/cmd/canary
```

## Usage

The first and only argument to canary is the url to check. The environment variable
SAMPLE_INTERVAL defines the interval (in seconds) to check the url. When the variable is
unset, the default of 1 second is used. 

To run canary with the default sampling interval:

```sh
$ canary http://www.canary.io
2014-12-28T14:44:32Z http://www.canary.io 200 96 true
2014-12-28T14:44:33Z http://www.canary.io 200 92 true
2014-12-28T14:44:34Z http://www.canary.io 200 89 true
2014-12-28T14:44:35Z http://www.canary.io 200 124 true
^C
```

To run canary with a 5 second sampling interval:

```sh
$ SAMPLE_INTERVAL=5 canary http://www.canary.io
2015-02-18T00:23:45-05:00 http://www.canary.io 200 78.111876 true 
2015-02-18T00:23:50-05:00 http://www.canary.io 200 72.897346 true 
2015-02-18T00:23:55-05:00 http://www.canary.io 200 60.863369 true 
2015-02-18T00:24:00-05:00 http://www.canary.io 200 69.095778 true 
^C^C
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
