#!/usr/bin/env sh

install_from_build() {
	if [ ! -e ./wut ]; then
        echo "Binary is not built, please run \"make\" first"
        exit 1
    fi
	cp ./wut /bin/wut
	echo "Install from build successful"
}

install_from_download() {
	# TODO: Implement installation from download
	echo "TODO: Implement installation from download"
}

if [ "$(id -u)" -ne 0 ]; then
	echo "Install failed: Install script not run as root."
	exit 1
fi

# Go to root of script
cd "${0%/*}" || exit 1

# Determine if inside repository or run from the internet
if [ -e ./wut.go ]; then
    install_from_build
else
    install_from_download
fi

exit 0
