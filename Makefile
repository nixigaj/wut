.PHONY: build build-debug run install clean github-release-build

default: build-debug

WUT_DIRTY_VERSION := $(shell git describe --tags --always)

ifndef WUT_VERSION
WUT_VERSION := $(WUT_DIRTY_VERSION)
endif

build:
	@go build -ldflags="-s -w -X main.wutVersion=$(WUT_VERSION)" -o wut

build-debug:
	@go build -ldflags="-X main.wutVersion=$(WUT_VERSION)" -o wut

run:
	@./wut

install:
	@if [ $$(id -u) -ne 0 ]; then \
		echo "Install failed: install not run as root." && exit 1; \
	fi
	@if [ ! -e ./wut ]; then \
		echo "Install failed: binary is not built, please run \'make build\' first." && exit 1; \
	fi
	cp ./wut /usr/local/bin/wut
	@echo "Install from build successful"

clean:
	@rm -rf wut wut.exe release

github-release-build:
	@if [ "$(WUT_VERSION)" = "$(WUT_DIRTY_VERSION)" ]; then \
		echo "WUT_VERSION environment variable needs to be set to a specific version" && exit 1; \
	fi

	@mkdir -p release

	@env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.wutVersion=$(WUT_VERSION)" -o wut
	@tar -czf release/wut-darwin-amd64.tar.gz wut
	@env GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.wutVersion=$(WUT_VERSION)" -o wut
	@tar -czf release/wut-darwin-arm64.tar.gz wut

	@env GOOS=freebsd GOARCH="386" go build -ldflags="-s -w -X main.wutVersion=$(WUT_VERSION)" -o wut
	@tar -czf release/wut-freebsd-386.tar.gz wut
	@env GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w -X main.wutVersion=$(WUT_VERSION)" -o wut
	@tar -czf release/wut-freebsd-amd64.tar.gz wut
	@env GOOS=freebsd GOARCH=arm64 go build -ldflags="-s -w -X main.wutVersion=$(WUT_VERSION)" -o wut
	@tar -czf release/wut-freebsd-arm64.tar.gz wut
	@env GOOS=freebsd GOARCH=arm go build -ldflags="-s -w -X main.wutVersion=$(WUT_VERSION)" -o wut
	@tar -czf release/wut-freebsd-arm.tar.gz wut

	@env GOOS=linux GOARCH="386" go build -ldflags="-s -w -X main.wutVersion=$(WUT_VERSION)" -o wut
	@tar -czf release/wut-linux-386.tar.gz wut
	@env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.wutVersion=$(WUT_VERSION)" -o wut
	@tar -czf release/wut-linux-amd64.tar.gz wut
	@env GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X main.wutVersion=$(WUT_VERSION)" -o wut
	@tar -czf release/wut-linux-arm64.tar.gz wut
	@env GOOS=linux GOARCH=arm go build -ldflags="-s -w -X main.wutVersion=$(WUT_VERSION)" -o wut
	@tar -czf release/wut-linux-arm.tar.gz wut

	@env GOOS=windows GOARCH="386" go build -ldflags="-s -w -X main.wutVersion=$(WUT_VERSION)" -o wut.exe
	@zip -q release/wut-windows-386.zip wut.exe
	@env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.wutVersion=$(WUT_VERSION)" -o wut.exe
	@zip -q release/wut-windows-amd64.zip wut.exe
	@env GOOS=windows GOARCH=arm64 go build -ldflags="-s -w -X main.wutVersion=$(WUT_VERSION)" -o wut.exe
	@zip -q release/wut-windows-arm64.zip wut.exe
	@env GOOS=windows GOARCH=arm go build -ldflags="-s -w -X main.wutVersion=$(WUT_VERSION)" -o wut.exe
	@zip -q release/wut-windows-arm.zip wut.exe

	@rm wut wut.exe
