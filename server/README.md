# `what` server
This is a super simple HTTP server made for one thing, and one thing only: echoing your public IP-address as fast as possible. It is written in Rust using [may_minihttp](https://github.com/Xudong-Huang/may_minihttp). I use it for [ip.erix.dev:11313](http://ip.erix.dev:11313).

## [wrk](https://github.com/wg/wrk) benchmarks
### Nginx
```
Running 10s test @ http://ip.erix.dev/
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    17.59ms   39.34ms 328.41ms   94.63%
    Req/Sec   546.21    148.36     1.08k    85.42%
  10527 requests in 10.01s, 4.17MB read
Requests/sec:   1051.29
Transfer/sec:    426.06KB
```

### `what` server
```
Running 10s test @ http://ip.erix.dev:11313/
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    12.12ms   33.68ms 310.40ms   93.32%
    Req/Sec     1.25k   295.07     1.49k    84.21%
  24149 requests in 10.01s, 9.56MB read
Requests/sec:   2413.43
Transfer/sec:      0.96MB
```
