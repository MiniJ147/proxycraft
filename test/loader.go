package main

import (
	"log"
	"net"
)

const P = "127.0.0.1:25565" //proxy
const D = "127.0.0.1:25566" // host

const MS uint8 = 1
const NC uint8 = 2

// inital connect to proxy
// proxy will then tell us to when to or when not to kill or establish new connections
// ESTABLIH
// 1. establish new conncetion to minecraft server
// 2. request a new connection to proxy (1 more than inital)
// 2.5. the proxy will then detect that we already have an active connection from
// this ip so it is a client pipe. From here it will assign it to the queue of
// clients wait and now they can communicate
// 3. begin to pipe the new connections together

func main() {
	m, e := net.Dial("tcp", P)
	if e != nil {
		log.Fatal("failed", e)
	}

	buf := make([]byte, 1024)
	n, e := m.Read(buf)
	if e != nil {
		log.Fatal("failed read", e)
	}
	log.Println(string(buf[:n]))

	for {
		n, e := m.Read(buf)
		if e != nil {
			log.Fatal("failed read", e, n)
		}

		switch buf[0] {
		case NC:
			log.Println("new connection")
		default:
			log.Println("failed")
		}
	}
}
