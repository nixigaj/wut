.PHONY: build build-debug run clean install

default: build-debug

build:
	go build -ldflags="-s -w" -o what

build-debug:
	go build -o what

run:
	@./what

clean:
	@rm what

install:
	@./install.sh
