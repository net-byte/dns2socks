package cmd

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"github.com/thinkgos/go-socks5/ccsocks5"
)

var c = cache.New(30*time.Minute, 10*time.Minute)

//StartServer starts server
func StartServer(localAddr *string, socksAddr *string, dnsAddr *string, cached *bool) {
	addr, err := net.ResolveUDPAddr("udp", *localAddr)
	if nil != err {
		log.Fatalln("failed to resolve udp addr:", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if nil != err {
		log.Fatalln("failed to listen udp:", err)
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
		if *cached {
			answer := getAnswerFromCache(b)
			if answer != nil {
				respMsg := new(dns.Msg)
				respMsg.Unpack(b)
				respMsg.Answer = append(respMsg.Answer, *answer)
				respData, _ := respMsg.Pack()
				conn.WriteToUDP(respData, fromAddr)
				continue
			}
		}
		proxyConn, err := client.Dial("udp", *dnsAddr)
		if err != nil {
			log.Println(err)
			continue
		}
		proxyConn.Write(b)
		go func() {
			defer proxyConn.Close()
			buf := make([]byte, 4096)
			for {
				n, err := proxyConn.Read(buf)
				if err != nil || err == io.EOF {
					break
				}
				b := buf[:n]
				if *cached {
					setAnswerCache(b)
				}
				conn.WriteToUDP(b, fromAddr)
			}
		}()
	}
}

func getAnswerFromCache(data []byte) *dns.RR {
	msg := new(dns.Msg)
	msg.Unpack(data)
	q := msg.Question[0]
	if q.Qtype == dns.TypeA {
		if v, found := c.Get(q.Name); found {
			ret := v.(*dns.RR)
			return ret
		}
	}
	return nil
}

func setAnswerCache(data []byte) {
	msg := new(dns.Msg)
	msg.Unpack(data)
	q := msg.Question[0]
	if q.Qtype == dns.TypeA && len(msg.Answer) > 0 {
		c.Set(q.Name, &msg.Answer[len(msg.Answer)-1], cache.DefaultExpiration)
	}
}
