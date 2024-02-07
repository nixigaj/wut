<img align="left" alt="wut icon" src="icon.svg" height="128" style="margin-right: 1rem" />

# `wut`
A fast and simple command-line tool to check your public IP-address.
It can also double as a network connectivity checker.

## Features / Usage
- **Fast** — Quires multiple different APIs at once and returns the first response while discarding the others.
- **Simple** — Does one thing and does it well, with a minimal number of options, in a single source file, using only the Go standard library.
- **IPv4 and IPv6** — Prints both by default but can print only one with the `--ipv4`/`-4` and `--ipv6`/`-6` flags. The default behaviour can be changed to a specific version by setting the `WUT_DEFAULT_IP_VERSION` environment variable to `ipv4`/`4` or `ipv6`/`6`. To print both versions with the environment variable set use the `--both` or `-b` flag.
- **Short output** — Does a pretty print by default but can print only an address with no white-space using the `--short` or `-s` flag. This flag additionally requires the IP version be explicitly specified with the `--ipv4`/`-4` or `--ipv6`/`-6` flags or the `WUT_DEFAULT_IP_VERSION` environment variable.
- **Specify interface or local IP** — Use a specific interface name or local IP by passing the `--interface` or `-i` flag. If an interface name is specified the first IPv4 and/or IPv6 will be chosen local IP(s).
- **Custom API** — One or more custom HTTP API address(es) can be supplied with one or more `--api` or `-a` flag(s). This expects that the API responds with only the IP-address in plaintext, as only white-space is trimmed from the response. It should also support both IPv4 and IPv6 for full functionality. Unencrypted HTTP is used by default unless a protocol is specified, e.g. `https://`.
- **Custom timeout** — Use a custom API fetch timeout in seconds with the `--timeout` or `-t` flag. By default, the timeout is three seconds.
- **Verbose error output** — Print verbose error output with the `--verbose` flag.
- **Print version** — Print the program version with the `--version` or `-v` flag.
- **Print usage** — Print program usage instructions with the `--help` or `-h` flag.

## Install
Since `wut` is a standalone binary,
it can be downloaded for the applicable platform from the [releases page](https://github.com/nixigaj/wut/releases)
and run from anywhere.
To install it automatically to the command-line path, the command below can be run with elevated privileges.

Darwin (macOS) and Linux:
```shell
curl -sSL https://raw.githubusercontent.com/nixigaj/wut/master/install.sh | sh
```

Darwin (macOS) and Linux with sudo inserted:
```shell
curl -sSL https://raw.githubusercontent.com/nixigaj/wut/master/install.sh | sudo sh
```

FreeBSD:
```shell
fetch -qo - https://raw.githubusercontent.com/nixigaj/wut/master/install.sh | sh
```

FreeBSD with doas inserted:
```shell
fetch -qo - https://raw.githubusercontent.com/nixigaj/wut/master/install.sh | doas sh
```

Windows:
```powershell
powershell -ExecutionPolicy Unrestricted -Command "Invoke-RestMethod -Uri https://raw.githubusercontent.com/nixigaj/wut/master/install.ps1 | Invoke-Expression"
```

Prebuilt binaries are available for:

| OS             | `386` | `amd64` | `arm` | `arm64` |
|----------------|-------|---------|-------|---------|
| Darwin (macOS) |       | ✅       |       | ✅       |
| FreeBSD        | ✅     | ✅       | ✅     | ✅       |
| Linux          | ✅     | ✅       | ✅     | ✅       |
| Windows        | ✅     | ✅       | ✅     | ✅       |

If your platform is not in the table, you can try building it from source below.

### Build from source
#### Dependencies
- Go 1.16 or higher
- Git
- Make (not required for Windows)

#### Clone repository and enter it

```shell
git clone https://github.com/nixigaj/wut.git
cd wut
```

#### Build

```shell
make build
```

#### Install
Run this command with elevated privileges:

```shell
make install
```

#### Windows
On Windows `make` can be replaced with `.\make.bat` in the commands.

## APIs

By default `wut` uses:

- [api64.ipify.org](https://api64.ipify.org)
- [icanhazip.com](https://icanhazip.com) ([this is usually the first one to respond](https://blog.apnic.net/2021/06/17/how-a-small-free-ip-tool-survived/))
- [ifconfig.me/ip](https://ifconfig.me/ip)
- [ip.erix.dev:11313](https://ip.erix.dev:11313) ([my own service](https://github.com/nixigaj/wut-server) in Sweden) (HTTP/2 only)
- [ipecho.net/plain](https://ipecho.net/plain)

### Roll your own API with Nginx
Use this directive and make sure that Nginx is not behind a reverse HTTP proxy:
```
return 200 "$remote_addr";
```

If you are feeling brave, you can also try the [Rust-based server](https://github.com/nixigaj/wut-server) that I use for [ip.erix.dev:11313](https://ip.erix.dev:11313).

## License
All files in this repository are licensed under the [MIT License](LICENSE).

The icon is a reference to the [Confused Nick Young / Swaggy P](https://knowyourmeme.com/memes/confused-nick-young-swaggy-p) meme.
