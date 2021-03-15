# dns2socks  

DNS proxy server over socks5(udp)

# Usage
## Cmd
```
Usage of ./bin/dns2socks:
  -c cache dns type a (default true)
  -d string
    	remote dns server address (default "8.8.8.8:53")
  -l string
    	local dns server address (default "127.0.0.1:53")
  -s string
    	socks5(udp) proxy address (default "127.0.0.1:1080")

```
## Docker  
```
docker run  -d --restart=always  --net=host --name dns2socks -p 53:53/udp netbyte/dns2socks -l=127.0.0.1:53 -s=127.0.0.1:1080 -d=8.8.8.8:53
```
