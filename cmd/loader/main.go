package main

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/minij147/proxycraft/pkg/consts"
)

var IP_DEST = "192.168.1.145:25566"

func NewClient(prox net.Conn, ip string) {

	dest, e := net.Dial("tcp", IP_DEST)
	if e != nil {
		log.Println("failed to connect to dest", e)
		return
	}
	defer dest.Close()

	log.Println(ip)
	pipe, e := net.Dial("tcp", consts.IP_PROXY_CONN)
	if e != nil {
		log.Println("failed to create pipe", e)
		return
	}
	defer pipe.Close()

	_, e = pipe.Write([]byte{consts.FLAG_CONN_OK})
	if e != nil {
		log.Println("failed to write throuhgh pipe")
		return
	}

	go func() {
		_, e := io.Copy(dest, pipe)
		if e != nil {
			log.Println(e)
		}
	}()

	_, e = io.Copy(pipe, dest)
	if e != nil {
		log.Println(e)
	}
}

func main() {
	log.Println("starting loader...")
	if len(os.Args) > 1 {
		IP_DEST = os.Args[1]
	}

	p, e := net.Dial("tcp", consts.IP_PROXY_CONN)
	if e != nil {
		log.Println("failed to connection", e)
	}

	_, e = p.Write([]byte{consts.FLAG_INIT})
	if e != nil {
		log.Println("failed to write", e)
	}

	buf := make([]byte, consts.PACKET_SIZE)

	n, e := p.Read(buf)
	if e != nil {
		log.Fatal("failed to get response")
	}

	if buf[0] != consts.FLAG_INIT_OK {
		log.Fatal("failed to init")
	}

	log.Println("initlized with server", string(buf[1:n]))
	log.Println(p.RemoteAddr().String())

	for {
		n, e := p.Read(buf)
		if e != nil {
			log.Fatal("failed to read to buf", e)
		}

		switch buf[0] {
		case consts.FLAG_CONN_NEW:
			go NewClient(p, string(buf[1:n]))
		case consts.FLAG_POLL:
			log.Println("polling answer")
		default:
			log.Println("unkown flag sent")
		}
	}
}
