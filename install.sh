#!/bin/sh

if [ "$(id -u)" -ne 0 ]; then
	echo "Install failed: install script not run as root."
	exit 1
fi

download_curl() {
	curl -LJ "https://github.com/nixigaj/wut/releases/latest/download/wut-${target}.tar.gz" -o "./wut-${target}.tar.gz"
}

download_freebsd_fetch() {
	fetch -o "./wut-${target}.tar.gz" "https://github.com/nixigaj/wut/releases/latest/download/wut-${target}.tar.gz"
}

case $(uname -ms) in
'Darwin x86_64')
	echo "Identified Darwin (macOS) and amd64 CPU"
	target=darwin-amd64 && download_curl
	;;
'Darwin arm64')
	echo "Identified Darwin (macOS) and arm64 CPU"
	target=darwin-arm64 && download_curl
	;;
'Linux i386' | 'Linux i686')
	echo "Identified Linux and 386 CPU"
	target=linux-386 && download_curl
	;;
'Linux amd64' | 'Linux x86_64')
	echo "Identified Linux and amd64 CPU"
	target=linux-amd64 && download_curl
	;;
'Linux aarch32' | 'Linux arm' | 'Linux arm32' | 'Linux armv6l' | 'Linux armv7l')
	echo "Identified Linux and arm CPU"
	target=linux-arm && download_curl
	;;
'Linux aarch64' | 'Linux arm64' | 'Linux armv8l')
	echo "Identified Linux and arm64 CPU"
	target=linux-arm64 && download_curl
	;;
'FreeBSD i386' | 'FreeBSD i686')
	echo "Identified FreeBSD and 386 CPU"
	target=freebsd-386 && download_freebsd_fetch
	;;
'FreeBSD amd64' | 'FreeBSD x86_64')
	echo "Identified FreeBSD and amd64 CPU"
	target=freebsd-amd64 && download_freebsd_fetch
	;;
'FreeBSD aarch32' | 'FreeBSD arm' | 'FreeBSD arm32' | 'FreeBSD armv6l' | 'FreeBSD armv7l')
	echo "Identified FreeBSD and arm CPU"
	target=freebsd-arm && download_freebsd_fetch
	;;
'FreeBSD aarch64' | 'FreeBSD arm64' | 'FreeBSD armv8l')
	echo "Identified FreeBSD and arm64 CPU"
	target=freebsd-arm64 && download_freebsd_fetch
	;;
*)
	echo "Install failed: platform $(uname -ms) is not supported."
	exit 1
	;;
esac

if [ ! -e "./wut-${target}.tar.gz" ]; then
	echo "Install failed: file not downloaded."
	exit 1
fi

tar -xzf "./wut-${target}.tar.gz" -C /usr/local/bin
rm -f "./wut-${target}.tar.gz"

echo "Installation successful."
exit 0
