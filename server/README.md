# `what` server
This is a super simple HTTP(/2) server written in Rust,
made for one thing, and one thing only: securely echoing your public IP-address as fast as possible.
I use it for [ip.erix.dev:11313](http://ip.erix.dev:11313).

## Non-scientific comparison with Nginx

The tool used for the benchmark is [wrk](https://github.com/wg/wrk).
The command used is `wrk -c 100 -t 8 <url>` over a 100 mb/s client connection.

### Nginx
```
Running 10s test @ https://ip.erix.dev
  8 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    19.14ms   17.07ms 363.55ms   98.11%
    Req/Sec   682.44     96.19     0.92k    88.79%
  54005 requests in 10.01s, 22.92MB read
Requests/sec:   5396.95
Transfer/sec:      2.29MB
```

### `what` server
```
Running 10s test @ https://ip.erix.dev:11313
  8 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    18.20ms   22.64ms 389.65ms   98.11%
    Req/Sec   764.16     82.61     1.14k    91.75%
  60168 requests in 10.01s, 7.52MB read
Requests/sec:   6013.30
Transfer/sec:    769.28KB
```

That is about than 11% faster.

`what` server has a lower transfer/sec, even though it has a higher requests/sec, because it sends less stuff in the response headers per request. By artificially inflating the response size with junk to match Nginx, you can still have higher requests/sec:
```
Running 10s test @ https://ip.erix.dev:11313
  8 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    17.02ms   11.81ms 311.18ms   98.04%
    Req/Sec   738.87     78.57     0.91k    86.87%
  58430 requests in 10.01s, 23.29MB read
Requests/sec:   5838.42
Transfer/sec:      2.33MB
```
