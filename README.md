<img align="left" alt="" src="icon.svg" height="128" style="margin-right: 1rem"/>

# `what`
A fast and simple command-line tool to check your public IP-address.
It can also double a network connectivity checker.

## Features / Usage
- **Fast** — Quires multiple different APIs at once and returns the first response while discarding the others.
- **Simple** — Does one thing and does it well, with a minimal number of options.
- **IPv4 and IPv6** — Does a pretty print of both by default but can print only the address with no white-space using the `-4` and `-6` flags.
- **Specify interface or gateway** — Use a specific interface name or gateway IP by passing the `--interface` or `-i` flag. If an interface name is specified the first IPv4 and/or IPv6 will be chosen as gateway.
- **Custom API** — One or more custom HTTP API address(es) can be supplied with one or more `--api` or `-a` flag(s). This expects that the API only responds with the IP-address in plaintext, as only white-space is trimmed from the response.
- **Print version** — Print the program version with the `--version` or `-v` flag(s).
- **Print usage** — Print program usage with the `--help` or `-h` flag(s).

### Planned
- **Optional curl backend** — Use [curl](https://curl.se) as backend for fetching the API(s) with the `--curl` or `-c` flag. This requires `curl` to be in the path.

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
Since `what` is a standalone binary,
it can be downloaded for the applicable platform from the [releases page](https://github.com/nixigaj/what/releases)
and run from anywhere.
To install it automatically to the command-line path, the command below can be run with elevated privileges.

Darwin (macOS), FreeBSD, and Linux:
```shell
curl -sSL https://raw.githubusercontent.com/nixigaj/zerve/master/install.sh | sh
```

Windows:
```powershell
curl -s https://raw.githubusercontent.com/nixigaj/zerve/master/install.bat | cmd
```

## APIs

By default `what` uses:

- [ip.erix.dev](https://ip.erix.dev) (my own service)
- [icanhazip.com](https://icanhazip.com)
- [ipecho.net/plain](https://ipecho.net/plain)
- [ifconfig.me/ip](https://ifconfig.me/ip)
- [api64.ipify.org](https://api64.ipify.org)

### Roll your own API with Nginx
Simply use this directive and make sure that Nginx is not behind some type of reverse proxy:
```
location / {
	default_type text/plain;
	return 200 "$remote_addr";
}
```

## License
All files in this repository are licensed under the [MIT License](LICENSE).

The icon is a reference to the [Confused Nick Young / Swaggy P](https://knowyourmeme.com/memes/confused-nick-young-swaggy-p) meme.
