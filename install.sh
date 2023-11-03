#!/bin/sh

if [ "$(id -u)" -ne 0 ]; then
	echo "Install failed: Install script not run as root."
	exit 1
fi

# TODO: Implement installation from download
echo "TODO: Implement installation from download"

exit 0
