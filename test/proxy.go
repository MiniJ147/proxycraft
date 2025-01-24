package main

import (
	"log"
	"net"
	"strings"
)

const H = "127.0.0.1:25565" // our address
var m map[string]net.Conn = make(map[string]net.Conn)
var x []net.Conn = make([]net.Conn,0 )

const MS uint8 = 1
const NC uint8 = 2

func hl(c net.Conn) {
	c.Write([]byte{MS})

	for {

	}

}

func hc(c net.Conn, s net.Conn) {
	defer c.Close()
	c.Write([]byte("hello"))

	s.Write([]byte{NC})

}

func main() {
	l, e := net.Listen("tcp", H)
	if e != nil {
		log.Fatal(e, "failed listen")
	}

	for {
		c, e := l.Accept()
		if e != nil {
			log.Println("failed conn", e)
			continue
		}

		ip := strings.Split(c.RemoteAddr().String(), ":")[0]
		v, ok := m[ip]
		if ok {
			log.Println("client")
			go hc(c, v)
		} else {
			log.Println("new connection")
			m[ip] = c
			go hl(c)
		}
	}
}
