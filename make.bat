@echo off

if "%1" == "setup" goto :setup
if "%1" == "build" goto :build
if "%1" == "build-debug" goto :build
if "%1" == "run" goto :run
if "%1" == "clean" goto :clean
if "%1" == "install" goto :install

REM Default target
if "%1" == "" goto :build-debug

echo Invalid target: %1
echo Usage: .\make.bat [setup^|build^|build-debug^|run^|clean^|install]
goto :eof

:setup
	go mod download
	goto :eof

:build
	go build -o what.exe
	goto :eof

:build-debug
	go build -ldflags="-s -w" -o what.exe
	goto :eof

:run
	.\what.exe
	goto :eof

:clean
	del what.exe
	goto :eof

:install
	.\install.bat
	goto :eof
