.PHONY: setup build run clean install

default: build-debug

setup:
	go mod download

build:
	go build -ldflags="-s -w" -o what *.go

build-debug:
	go build -o what *.go

run:
	@./what

clean:
	@rm what

install:
	@./install.sh
