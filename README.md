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
$ canary -u http://www.canary.io
2014-11-12T22:54:33Z resh.local http://www.canary.io 199.27.79.133 0 200 0.343872
2014-11-12T22:54:34Z resh.local http://www.canary.io 199.27.79.133 0 200 0.142661
2014-11-12T22:54:35Z resh.local http://www.canary.io 199.27.79.133 0 200 0.141929
2014-11-12T22:54:36Z resh.local http://www.canary.io 199.27.79.133 0 200 0.080171
#...
```

Do the same, but emit [logfmt](https://engineering.heroku.com/blogs/2014-09-05-hutils-explore-your-structured-data-logs):

```sh
$ canary -u http://www.canary.io -o logfmt
t=2014-11-12T22:56:44Z source=resh.local url=http://www.canary.io ip=23.235.47.133 curl_status=0 ht
tp_status=200 namelookup_time=0.468537 connect_time=0.531856 start_transfer_time=0.672994 total_tim
e=0.743690                                                                                        
t=2014-11-12T22:56:45Z source=resh.local url=http://www.canary.io ip=23.235.47.133 curl_status=0 ht
tp_status=200 namelookup_time=0.000036 connect_time=0.000036 start_transfer_time=0.075379 total_tim
e=0.144114                                                                                        
t=2014-11-12T22:56:46Z source=resh.local url=http://www.canary.io ip=23.235.47.133 curl_status=0 ht
tp_status=200 namelookup_time=0.000032 connect_time=0.000032 start_transfer_time=0.071334 total_tim
e=0.134759                                                                                        
t=2014-11-12T22:56:47Z source=resh.local url=http://www.canary.io ip=23.235.47.133 curl_status=0 ht
tp_status=200 namelookup_time=0.000048 connect_time=0.000048 start_transfer_time=0.065151 total_tim
e=0.069570                                                                                        
t=2014-11-12T22:56:48Z source=resh.local url=http://www.canary.io ip=23.235.47.133 curl_status=0 ht
tp_status=200 namelookup_time=0.000059 connect_time=0.000060 start_transfer_time=0.063251 total_tim
e=0.066277
# ...
```
