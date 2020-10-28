#!/bin/bash

build_dir=`pwd`

echo $(date +"%H:%M:%S")  "start darwin compile"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o build/sync-mysql-darwin
echo $(date +"%H:%M:%S")  "finish darwin compile"

echo $(date +"%H:%M:%S")  "start linux compile"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o build/sync-mysql-linux
echo $(date +"%H:%M:%S")  "finish linux compile"

echo $(date +"%H:%M:%S")  "start windows compile"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o build/sync-mysql-windows.exe
echo $(date +"%H:%M:%S")  "finish windows compile"

