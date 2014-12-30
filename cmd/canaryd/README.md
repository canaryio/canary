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

## Manifest

A manifest is a simple JSON document describing the sites to be monitored.  An example:

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
      "name": "github"
    }
  ]
}
```

## Usage

```sh
$ MANIFEST_URL=http://www.canary.io/manifest.json canaryd
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
