package main

import (
	"fmt"
	"io"
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

// TODO: add dynamic changes to ports and ips and pull it out of hard coded constants
// TODO: add error messages  for client

func main() {
	log.Println("starting the loader...")

	proxy, err := net.DialTimeout("tcp", consts.IP_PROXY, TIMEOUT_PROXY)
	if err != nil {
		log.Fatal("failed  to connect  to proxy on", consts.IP_PROXY)
	}

	dest, err := net.DialTimeout("tcp", IP_DEST, TIMEOUT_DEST)
	if err != nil {
		log.Fatal("failed to connect to dest on ip", IP_DEST)
	}

	buf := make([]byte, consts.PACKET_SIZE)

	_, err = proxy.Write([]byte{consts.FLAG_CREATE})
	if err != nil {
		log.Fatal("failed create packet to proxy", err)
	}

	n, err := proxy.Read(buf)
	if err != nil || n == 0 {
		log.Fatal("failed to read buffer", err)
	}

	if buf[0] == consts.FLAG_FAIL {
		log.Fatal("failed to create server")
	}

	log.Println(string(buf[1:n]))
	proxy.Write(packets.Create(consts.FLAG_DATA, 1, consts.PACKET_EMPTY))

	//NOTE: TEST CODE (only 1 users at a time and disconnects not supported)
	log.Println(io.Copy(dest, proxy))

	log.Println("Done, press any key to close...")
	fmt.Scanf("")
}
