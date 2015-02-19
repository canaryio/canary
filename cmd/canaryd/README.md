canaryd
======

a mutli-site monitoring tool.

## Installation

```sh
$ go get github.com/canaryio/canary/cmd/canaryd
```

## Configuration

`canaryd` is configured via environment variables:

* `MANIFEST_URL` - ref to a JSON document describing what needs to be monitored
* `PUBLISHERS` - an explicit list of pubilshers to enable, defaulting to `stdout`
* `DEFAULT_SAMPLE_INTERVAL` - interval rate (in seconds) for targets without a defined interval value, defaults to 1 second.

## Manifest

A manifest is a simple JSON document describing the sites to be monitored.  You must create such a document and host it somewhere so that it is accessible to `canaryd`.

Within the manifest, targets are defined as a json object with the required keys 'url' and 'name'. 'interval' is optional, and will define the interval rate in seconds to check the specific url, overriding the default interval settings in canaryd

An example manifest:

```js
{
  "targets": [
    {
      "url": "http://www.canary.io",
      "name": "canary"
    },
    {
      "url": "https://www.simple.com/",
      "name": "simple"
    },
    {
      "url": "https://www.heroku.com/",
      "name": "heroku"
    },
    {
      "url": "https://github.com",
      "name": "github",
      "interval": 60
    }
  ]
}
```

## Publishers

`canaryd` supports a number of configurable publishers.

### `stdout`

The default publisher.  Writes all measurements to `STDOUT` as they happen.

If `PUBLISHERS` are not set, this is activated by default.

To explicitly activate, set `PUBLISHERS=stdout`.

Example:

```sh
$ PUBLISHERS=stdout MANIFEST_URL=http://www.canary.io/manifest.json canaryd
2014/12/27 15:20:09 http://www.canary.io 200 128 true
2014/12/27 15:20:09 https://www.simple.com/ 200 252 true
2014/12/27 15:20:09 https://github.com 200 384 true
2014/12/27 15:20:09 https://www.heroku.com/ 200 413 true
2014/12/27 15:20:10 https://www.simple.com/ 200 76 true
2014/12/27 15:20:10 http://www.canary.io 200 94 true
2014/12/27 15:20:10 https://github.com 200 306 true
2014/12/27 15:20:10 https://www.heroku.com/ 200 306 true
^C
```

### `librato`

Writes all measurements to your [Librato](https://www.librato.com/) account at 5 second intervals.

To activate, set `PUBLISHERS=librato`.

For configuration purposes, the Librato publisher expects the following environment variables to bet set:

| Variable | Required | Description |
| -------- | -------- | ----------- |
| `LIBRATO_USER` | Yes | Librato API user |
| `LIBRATO_TOKEN` | Yes | Librato API token |
| `SOURCE` | No | source name to use in metrics, defaults to [`os.Hostname`](http://golang.org/pkg/os/#Hostname) |

The following metrics are produced:

| Metric | Description |
| ------ | ----------- |
| `canary.{NAME}.latency` | the time it took to complete the `GET` request |
| `canary.{NAME}.errors` | a count of samples that included an error |
| `canary.{NAME}.errors.http` | a count of samples that contained HTTP status codes outside of the 3xx range |
| `canary.{NAME}.errors.sampler` | a count of samples that indicated transport-level error such as a timeout or connection failure |

An example invocation:

```sh
$ PUBLISHERS=librato LIBRATO_USER=michael.gorsuch@gmail.com LIBRATO_TOKEN=REDACTED MANIFEST_URL=http://www.canary.io/manifest.json canaryd
#...
```