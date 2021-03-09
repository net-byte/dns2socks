package cmd

import (
	"io"
	"log"
	"net"

	"github.com/miekg/dns"
	"github.com/thinkgos/go-socks5/ccsocks5"
)

//StartServer acts start a dns proxy server
func StartServer(localAddr *string, socksAddr *string, dnsAddr *string) {
	addr, err := net.ResolveUDPAddr("udp", *localAddr)
	if nil != err {
		log.Fatalln("Unable to get UDP socket:", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if nil != err {
		log.Fatalln("Unable to listen on UDP socket:", err)
	}
	defer conn.Close()
	log.Printf("dns2socks started on %v", *localAddr)
	client := ccsocks5.NewClient(*socksAddr)
	defer client.Close()
	buf := make([]byte, 4096)
	for {
		n, fromAddr, err := conn.ReadFromUDP(buf)
		if err != nil || err == io.EOF {
			continue
		}
		b := buf[:n]
		printQuestion(b)
		proxyConn, err := client.Dial("udp", *dnsAddr)
		if err != nil {
			log.Println(err)
		}
		defer proxyConn.Close()
		proxyConn.Write(b)
		go func() {
			buf := make([]byte, 4096)
			for {
				n, err := proxyConn.Read(buf)
				if err != nil || err == io.EOF {
					break
				}
				b := buf[:n]
				printAnswer(b)
				conn.WriteToUDP(b, fromAddr)
			}
		}()
	}
}

func printQuestion(data []byte) {
	msg := new(dns.Msg)
	msg.Unpack(data)
	log.Printf("dns question:%v", msg.Question)
}

func printAnswer(data []byte) {
	msg := new(dns.Msg)
	msg.Unpack(data)
	log.Printf("dns answer:%v", msg.Answer)
}
