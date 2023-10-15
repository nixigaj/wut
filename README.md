<img align="left" alt="" src="icon.svg" height="128" style="margin-right: 1rem"/>

# `what`
A fast and simple command-line tool to check your public IP-address. It can also double a network connectivity checker.

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

## Install


Since `what` is a standalone binary, it can be downloaded for the applicable platform from the [releases page](https://github.com/nixigaj/what/releases) and run from anywhere. To install it automatically to the command-line path, the command below can be run with elevated privileges.

Darwin (macOS), FreeBSD, and Linux:

```shell
curl -sSL https://raw.githubusercontent.com/nixigaj/zerve/master/install.sh | sh
```

Windows:

```powershell
curl -s https://raw.githubusercontent.com/nixigaj/zerve/master/install.bat | cmd
```

## License
All files in this repository are licensed under the [MIT License](LICENSE).
