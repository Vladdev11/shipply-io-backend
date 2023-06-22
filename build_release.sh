#!/bin/bash
export PATH=$PATH:~/go/bin
go generate
env GOOS=linux GOARCH=amd64 go build -o "bin_release/shipply-client-backend" -trimpath -ldflags "-w -s"