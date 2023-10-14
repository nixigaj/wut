# `what`
A fast and simple command-line tool to check your public IP address. It can also double a network connectivity checker.

## Features / Usage
- **Fast** — Quires multiple different APIs at once and returns the first response while discarding the others.
- **Simple** — Does one thing and does it well, with a minimal amount of options.
- **IPv4 and IPv6** — Does a pretty print of both by default but can print only the address with the `-4` and `-6` flags.
- **Specify interface or gateway** — Use a specific interface name or gateway IP by passing the `--interface` or `-i` flag.
- **Custom API** — One or more custom HTTP API address(es) can be supplied with one or more `--api` or `-a` flag(s).

### Planned
- **Optional curl backend** — Use curl as backend for fetching the API(s) with the `--curl` or `-c` flag. This requires `curl` to be in the path.

## Platform support
- Darwin (macOS)
- FreeBSD
- Linux
- Windows

### Planned
- DragonFly BSD
- NetBSD
- OpenBSD
- Solaris
