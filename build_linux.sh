#!/bin/bash

build_dir=`pwd`

echo $(date +"%H:%M:%S")  "start compile"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o sync-mysql
echo $(date +"%H:%M:%S")  "finish compile"
