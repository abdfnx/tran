# get latest release
$release_url = "https://api.github.com/repos/abdfnx/tran/releases"
$tag = (Invoke-WebRequest -Uri $release_url -UseBasicParsing | ConvertFrom-Json)[0].tag_name
$loc = "$HOME\AppData\Local\tran"
$url = ""
$arch = $env:PROCESSOR_ARCHITECTURE
$releases_api_url = "https://github.com/abdfnx/tran/releases/download/$tag/tran_windows_${tag}"

if ($arch -eq "AMD64") {
    $url = "${releases_api_url}_amd64.zip"
} elseif ($arch -eq "x86") {
    $url = "${releases_api_url}_386.zip"
} elseif ($arch -eq "arm") {
    $url = "${releases_api_url}_arm.zip"
} elseif ($arch -eq "arm64") {
    $url = "${releases_api_url}_arm64.zip"
}

if (Test-Path -path $loc) {
    Remove-Item $loc -Recurse -Force
}

Write-Host "Installing tran version $tag" -ForegroundColor DarkCyan

Invoke-WebRequest $url -outfile tran_windows.zip

Expand-Archive tran_windows.zip

New-Item -ItemType "directory" -Path $loc

Move-Item -Path tran_windows\bin -Destination $loc

Remove-Item tran_windows* -Recurse -Force

[System.Environment]::SetEnvironmentVariable("Path", $Env:Path + ";$loc\bin", [System.EnvironmentVariableTarget]::User)

if (Test-Path -path $loc) {
    Write-Host "Thanks for installing Tran! Now Refresh your powershell" -ForegroundColor DarkGreen
    Write-Host "If this is your first time using the CLI, be sure to run 'tran --help' first." -ForegroundColor DarkGreen
} else {
    Write-Host "Download failed" -ForegroundColor Red
    Write-Host "Please try again later" -ForegroundColor Red
}
