.PHONY: setup build run clean install

default: build

setup:
	go mod download

build:
	go build -o what *.go

run:
	@./what

clean:
	@rm what

install: build
	@./install.sh
