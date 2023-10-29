@echo off

if "%1" == "build" goto :build
if "%1" == "build-debug" goto :build
if "%1" == "run" goto :run
if "%1" == "clean" goto :clean
if "%1" == "install" goto :install

REM Default target
if "%1" == "" goto :build-debug

echo Invalid target: %1
echo Usage: .\make.bat [build^|build-debug^|run^|clean^|install]
goto :eof

:build
	go build -ldflags="-s -w" -o wut.exe
	goto :eof

:build-debug
	go build -o wut.exe
	goto :eof

:run
	.\wut.exe
	goto :eof

:clean
	del wut.exe
	goto :eof

:install
	REM TODO: Write Windows installation script.
	echo TODO: Write Windows installation script
	goto :eof
