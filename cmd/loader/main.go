package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/minij147/proxycraft/pkg/consts"
	"github.com/minij147/proxycraft/pkg/packets"
)

const (
	//NOTE:
	// proxy and dest will both be on port 25565 but shared different addresses
	// I will add custom ip and ports later

	IP_DEST = "127.0.0.1:25566" // minecraft

	TIMEOUT_PROXY = 30 * time.Second
	TIMEOUT_DEST  = 30 * time.Second
)

// shouldn't need to be sync (hopefully)
var clients map[uint32]net.Conn = make(map[uint32]net.Conn)

// TODO: add dynamic changes to ports and ips and pull it out of hard coded constants
// TODO: add error messages  for client

// Returns Proxy, Dest
func Init() (net.Conn, net.Conn) {
	proxy, err := net.DialTimeout("tcp", consts.IP_PROXY, TIMEOUT_PROXY)
	if err != nil {
		log.Fatal("failed  to connect  to proxy on", consts.IP_PROXY)
	}

	// dest, err := net.DialTimeout("tcp", IP_DEST, TIMEOUT_DEST)
	// if err != nil {
	// 	log.Fatal("failed to connect to dest on ip", IP_DEST)

	buf := make([]byte, consts.PACKET_SIZE_SIGNED)

	_, err = proxy.Write([]byte{consts.FLAG_CREATE})
	if err != nil {
		log.Fatal("failed create packet to proxy", err)
	}

	n, err := proxy.Read(buf)
	if err != nil || n == 0 {
		log.Fatal("failed to read buffer", err)
	}
	flag, _, data := packets.Read(n, buf)

	if flag == consts.FLAG_FAIL {
		log.Fatal("failed to create server")
	}

	log.Println(string(data))

	return proxy, nil
}

// is to sign the data from the server back up to proxy
// main loop will write data on our behalf
func clientSignLifetime(id uint32, src net.Conn, dest net.Conn) {
	buf := make([]byte, consts.PACKET_SIZE_RAW)
	for {
		n, err := src.Read(buf)
		if err != nil {
			// TODO: add fail safes
			log.Println(id, "failed read to server")
			break
		}

		packet := packets.Create(consts.FLAG_DATA, id, buf[:n])
		_, err = dest.Write(packet)
		if err != nil {
			log.Println(id, "failed to write to dest")
			break
		}
	}

	log.Println(id, "shutting down")
}

func main() {
	log.Println("starting the loader...")
	proxy, _ := Init()

	// all packets from here will be signed via the proxy
	buf := make([]byte, consts.PACKET_SIZE_SIGNED)
	for {
		n, err := proxy.Read(buf)
		if err != nil {
			log.Println("read from proxy error", err)
			break
		}

		flag, id, data := packets.Read(n, buf)
		switch flag {
		case consts.FLAG_DATA:
			log.Println(buf[:n], data)
			client, ok := clients[id]
			if !ok {
				//TODO: add removal later
				log.Println("could not find client with id of", id)
				break
			}

			_, err := client.Write(data)
			if err != nil {
				//TODO: add error handling later
				log.Println("failed to write packet to client", id)
				break
			}
		case consts.FLAG_CONNECTION_INCOMING:
			fail := func() {
				packet := packets.Create(consts.FLAG_CONNECTION_FAILED, id, consts.PACKET_EMPTY)
				proxy.Write(packet)
			}
			log.Println("connection incoming", id)

			_, ok := clients[id]
			if ok {
				log.Println("failed incoming connection client already has connection", id)
				fail()
				break
			}

			newConn, err := net.DialTimeout("tcp", IP_DEST, TIMEOUT_DEST)
			if err != nil {
				log.Println("failed to establish new connection for client", id, err)
				fail()
				break
			}

			clients[id] = newConn
			_, err = newConn.Write(data) // writing inital data
			if err != nil {
				log.Println("failed to write inital to confrim connection to server", err)
				fail()
				break
			}

			go clientSignLifetime(id, newConn, proxy)

			// add proper  solution later for proxy going down
			_, err = proxy.Write(packets.Create(consts.FLAG_CONNECTION_ACCEPTED, id, consts.PACKET_EMPTY))
			if err != nil {
				log.Println("failed to write to proxy")
				break
			}

			log.Println("connection made")

		default:
			log.Println("unkown flag", flag, id)
		}
	}

	log.Println("Done, press any key to close...")
	fmt.Scanf("")
}
