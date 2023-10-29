.PHONY: build build-debug run clean install

default: build-debug

build:
	go build -ldflags="-s -w" -o wut

build-debug:
	go build -o wut

run:
	@./wut

clean:
	@rm wut

install:
	@if [ $$(id -u) -ne 0 ]; then \
		echo "Install failed: Install script not run as root." && exit 1; \
	fi
	@if [ ! -e ./wut ]; then \
		echo "Binary is not built, please run \`make build\` first" && exit 1; \
	fi
	cp ./wut /usr/local/bin/wut
	@echo "Install from build successful"
