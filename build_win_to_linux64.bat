@echo off

set GOOS=linux
set GOARCH=amd64

echo "Build For sync-mysql ..."

go build -ldflags "-s -w" -o sync-mysql

echo "--------- Build For sync-mysql Success!"
