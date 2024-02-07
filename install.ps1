$wutDirectory = "${env:ProgramFiles}\wut"

function Test-IsAdministrator {
	return ([Security.Principal.WindowsPrincipal]`
	        [Security.Principal.WindowsIdentity]::GetCurrent()`
	).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator) -and $env:USERNAME -ne "WDAGUtilityAccount"
}

if (!(Test-IsAdministrator)) {
	Write-Host "Install failed: install script not run as administrator."
	exit 1
}

function Install-Wut-Remote {
	param (
		[string]$CpuArch
	)

	$downloadUrl = "https://github.com/nixigaj/wut/releases/latest/download/wut-windows-$CpuArch.zip"

	New-Item -ItemType Directory -Force -Path $wutDirectory | Out-Null
	Invoke-WebRequest $downloadUrl -OutFile "wut-windows-$CpuArch.zip"
	Expand-Archive -Force "wut-windows-$CpuArch.zip" -DestinationPath $wutDirectory

	if (Test-Path "wut-windows-$CpuArch.zip") {
		Remove-Item "wut-windows-$CpuArch.zip" | Out-Null
	}
}

function Install-Wut-Build {
	if (!(Test-Path "wut.exe")) {
		Write-Host "Install failed: executable is not built. Please run '.\make.bat build'"
		exit 1
	}

	New-Item -ItemType Directory -Force -Path $wutDirectory | Out-Null
	Copy-Item -Force -Path .\wut.exe -Destination $wutDirectory\wut.exe
}

if ($env:WUT_INSTALL_BUILD -eq "true") {
	Write-Host "Installing from local build"
	Install-Wut-Build
} elseif ($env:PROCESSOR_ARCHITECTURE -eq "x86") {
	Write-Host "Identified 386 CPU"
	Install-Wut-Remote -CpuArch "386"
} elseif ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") {
	Write-Host "Identified amd64 CPU"
	Install-Wut-Remote -CpuArch "amd64"
} elseif ($env:PROCESSOR_ARCHITECTURE -eq "ARM") {
	Write-Host "Identified arm CPU"
	Install-Wut-Remote -CpuArch "arm"
} elseif ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
	Write-Host "Identified arm64 CPU"
	Install-Wut-Remote -CpuArch "arm64"
} else {
	Write-Host "Install failed: Unsupported processor type: $($env:PROCESSOR_ARCHITECTURE)"
	exit 1
}

$newPath = $wutDirectory
$currentPath = [System.Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::Machine)

if ($currentPath -notlike "*$newPath*") {
	$newPath = "$currentPath;$newPath"

	[System.Environment]::SetEnvironmentVariable("Path", $newPath, [System.EnvironmentVariableTarget]::Machine)

	Write-Host "Added wut to path. Restart terminal for changes to take effect."
}

Write-Host "Installation successful."
exit 0
