#!/bin/bash

# Build for 64-bit macOS
GOOS=darwin GOARCH=amd64 go build -o wtmonitor-client-mac-64 main.go fileupload.go

checksum=$(shasum -a 256 wtmonitor-client-mac-64)

echo "macOS 64-bit ✅"
echo $checksum

# Build for 32-bit Windows
GOOS=windows GOARCH=386 go build -o wtmonitor-client-windows-32.exe main.go fileupload.go

checksum=$(shasum -a 256 wtmonitor-client-windows-32)

echo "Windows 32-bit ✅"
echo $checksum

# Build for 32-bit Windows
GOOS=windows GOARCH=amd64 go build -o wtmonitor-client-windows-64.exe main.go fileupload.go

checksum=$(shasum -a 256 wtmonitor-client-windows-64)

echo "Windows 64-bit ✅"
echo $checksum

echo "Build Successful 🙂"