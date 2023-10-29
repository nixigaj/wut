.PHONY: build build-debug run clean install github-release-build

default: build-debug

build:
	go build -ldflags="-s -w" -o wut

build-debug:
	go build -o wut

run:
	@./wut

clean:
	@rm -rf wut wut.exe release

install:
	@if [ $$(id -u) -ne 0 ]; then \
		echo "Install failed: Install script not run as root." && exit 1; \
	fi
	@if [ ! -e ./wut ]; then \
		echo "Binary is not built, please run \`make build\` first" && exit 1; \
	fi
	cp ./wut /usr/local/bin/wut
	@echo "Install from build successful"

github-release-build:
	@mkdir -p release

	@env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o wut
	@tar -czf release/wut-darwin-amd64.tar.gz wut
	@env GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o wut
	@tar -czf release/wut-darwin-arm64.tar.gz wut

	@env GOOS=freebsd GOARCH="386" go build -ldflags="-s -w" -o wut
	@tar -czf release/wut-freebsd-386.tar.gz wut
	@env GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o wut
	@tar -czf release/wut-freebsd-amd64.tar.gz wut
	@env GOOS=freebsd GOARCH=arm64 go build -ldflags="-s -w" -o wut
	@tar -czf release/wut-freebsd-arm64.tar.gz wut
	@env GOOS=freebsd GOARCH=arm go build -ldflags="-s -w" -o wut
	@tar -czf release/wut-freebsd-arm.tar.gz wut

	@env GOOS=linux GOARCH="386" go build -ldflags="-s -w" -o wut
	@tar -czf release/wut-linux-386.tar.gz wut
	@env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o wut
	@tar -czf release/wut-linux-amd64.tar.gz wut
	@env GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o wut
	@tar -czf release/wut-linux-arm64.tar.gz wut
	@env GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o wut
	@tar -czf release/wut-linux-arm.tar.gz wut

	@env GOOS=windows GOARCH="386" go build -ldflags="-s -w" -o wut.exe
	@zip -q release/wut-windows-386.zip wut.exe
	@env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o wut.exe
	@zip -q release/wut-windows-amd64.zip wut.exe
	@env GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o wut.exe
	@zip -q release/wut-windows-arm64.zip wut.exe
	@env GOOS=windows GOARCH=arm go build -ldflags="-s -w" -o wut.exe
	@zip -q release/wut-windows-arm.zip wut.exe

	@rm wut wut.exe
