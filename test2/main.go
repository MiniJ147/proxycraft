package main

import (
	"log"
	"net"
	"time"
)

const PACKET_HEADER = 1
const ADJ = 100

// const DEST_IP = "127.0.0.1:25567"   // example host ip
const MIDDLE_IP = "127.0.0.1:25565" // our proxy host
// we don't need minecraft ip

func main() {
	log.Println("hello world!")

	proxy, err := net.Listen("tcp", MIDDLE_IP)
	if err != nil {
		log.Fatal("step 1", err)
	}
	log.Println("started listing")

	server, err := proxy.Accept()
	if err != nil {
		log.Fatal("failed server connection", err)
	}
	log.Println("have server")

	client, err := proxy.Accept()
	if err != nil {
		log.Fatal("failed player connection", err)
	}
	log.Println("have client")

	// server to client [passed through proxy]
	go func() {
		buf := make([]byte, ADJ+1024+PACKET_HEADER)
		for {
			n, err := server.Read(buf)
			if err != nil {
				server.Close()
				log.Println("cleint -> server err", err)
				break
			}
			log.Println(buf[:n])

			_, err = client.Write(buf[:n])
			if err != nil {
				log.Println("error ", err)
				break
			}
		}
	}()

	// client to loader [passed through proxy]
	buf := make([]byte, PACKET_HEADER+1024)
	for {
		n, err := client.Read(buf[PACKET_HEADER:])
		if err != nil {
			log.Println("cleint -> server err", err)
			break
		}
		buf[0] = 200
		// buf[1] = 0
		// buf[2] = 0
		// buf[3] = 0
		// buf[4] = 1

		// log.Println(buf)
		x, err := server.Write(buf[:n+PACKET_HEADER])
		if err != nil {
			log.Println("error ", err)
			break
		}
		log.Println(x)
	}
	// _, err = io.Copy(server, client)
	// if err != nil {
	// 	log.Println("client -> server", err)
	// }

	log.Println("no longer doing copies")
	time.Sleep(20 * time.Second)
	log.Println("killing program")

}
