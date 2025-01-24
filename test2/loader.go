package main

import (
	"log"
	"net"
)

const PACKET_HEADER = 1
const PROXY_IP = "127.0.0.1:25565"
const MINECRAFT_IP = "192.168.1.145:25566"

func main() {
	proxy, err := net.Dial("tcp", PROXY_IP)
	if err != nil {
		log.Fatal("failed connection", err)
	}

	log.Println("connected to proxy")
	server, err := net.Dial("tcp", MINECRAFT_IP)
	if err != nil {
		log.Fatal("failed onncetion to minecraft", err)
	}
	log.Println("connected to minecraft")

	// go func() {
	// 	_, err := io.Copy(proxy, server)
	// 	if err != nil {
	// 		log.Println("server -> proxy done")
	// 	}
	// }()
	//
	// _, err = io.Copy(server, proxy)
	// if err != nil {
	// 	log.Println("proxy -> server done")

	// }
	// receviing from minecraft to middleman
	go func() {
		for {
			buf := make([]byte, 1024+PACKET_HEADER)
			n, err := server.Read(buf)
			if err != nil {
				server.Close()
				log.Println("cleint -> server err", err)
				break
			}
			// buf[0] = 200
			// buf[1] = 0
			// buf[2] = 0
			// buf[3] = 0
			// buf[4] = 2
			log.Println(buf[:n])

			_, err = proxy.Write(buf[:n])
			if err != nil {
				log.Println("error ", err)
				break
			}
		}
		//       buf :=
		// _, err := io.Copy(client, server)
		// if err != nil {
		// 	log.Println("server -> client", err)
		// }
	}()

	// middleman to minecraft
	buf := make([]byte, PACKET_HEADER+1024)
	for {
		n, err := proxy.Read(buf)
		if err != nil {
			server.Close()
			log.Println("cleint -> server err", err)
			break
		}
		log.Println(buf[0], buf)

		_, err = server.Write(buf[PACKET_HEADER:n])
		if err != nil {
			log.Println("error ", err)
			break
		}
	}
}
