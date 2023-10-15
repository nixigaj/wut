#!/usr/bin/env sh

if [ "$(id -u)" -ne 0 ]; then
	echo "Install failed: Install script not run as root."
	exit 1
fi

# Go to root of repository
cd "${0%/*}" || exit 1

cp ./what /bin/what

exit 0
