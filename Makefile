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
	@./install.sh
