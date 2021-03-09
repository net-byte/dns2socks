package main

import (
	"flag"

	"github.com/net-byte/dns2socks/cmd"
)

func main() {
	localAddr := flag.String("l", "127.0.0.1:53", "local dns server address")
	socksAddr := flag.String("s", "127.0.0.1:1080", "socks5(udp) proxy address")
	dnsAddr := flag.String("d", "8.8.8.8:53", "remote dns server address")
	flag.Parse()
	cmd.StartServer(localAddr, socksAddr, dnsAddr)
}
