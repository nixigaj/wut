@echo off

if "%1" == "build" goto :build
if "%1" == "build-debug" goto :build
if "%1" == "run" goto :run
if "%1" == "install" goto :install
if "%1" == "clean" goto :clean

REM Default target
if "%1" == "" goto :build-debug

echo Invalid target: %1
echo Usage: .\make.bat [build^|build-debug^|run^|install^|clean]
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

:install
	setlocal enabledelayedexpansion
	set WUT_INSTALL_BUILD=true
	powershell -ExecutionPolicy Unrestricted -File .\install.ps1
	goto :eof

:clean
	if exist wut.exe del wut.exe
	goto :eof
