# `what` server
This is a super simple HTTP(/2) server written in Rust,
made for one thing, and one thing only: echoing your public IP-address as fast as possible.
I use it for [ip.erix.dev:11313](http://ip.erix.dev:11313).

## Non-scientific comparison with Nginx

The tool used for the benchmark is [wrk](https://github.com/wg/wrk).

### Nginx
```
Running 10s test @ https://ip.erix.dev/
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    19.91ms   54.03ms 460.43ms   95.32%
    Req/Sec   544.61    143.72     1.35k    88.66%
  10532 requests in 10.00s, 4.47MB read
Requests/sec:   1052.92
Transfer/sec:    457.56KB
```

### `what` server
```
Running 10s test @ https://ip.erix.dev:11313/
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    10.22ms   29.15ms 268.85ms   94.84%
    Req/Sec     1.25k   288.35     1.51k    87.76%
  24411 requests in 10.00s, 3.35MB read
Requests/sec:   2440.28
Transfer/sec:    343.16KB
```

`what` server has a lower transfer/sec, even though it has a higher requests/sec, because it sends less stuff in the response headers per request.
