#!bin/bash
export GO111MODULE=on

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-s -w" -o ./bin/dns2socks-openwrt-amd64 ./main.go

echo "done!"
