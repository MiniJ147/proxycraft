package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

const PROXY_IP = "127.0.0.1:25565"
const DEST_IP = "127.0.0.1:25566"

const PROXY_TIMEOUT = 30 * time.Second
const DEST_TIMEOUT = 30 * time.Second

const FLAG_CREATE = 100
const FLAG_SUCCESS = 101
const FLAG_FAIL = 102

const PACKET_SIZE = 1024

func main() {
	proxy, err := net.DialTimeout("tcp", PROXY_IP, PROXY_TIMEOUT)
	if err != nil {
		log.Fatal("failed connection to proxy server, check status", err)
	}

	dest, err := net.DialTimeout("tcp", DEST_IP, DEST_TIMEOUT)
	if err != nil {
		log.Fatal("failed to connect to your server", err)
	}

	_, err = proxy.Write([]byte{FLAG_CREATE})
	if err != nil {
		log.Fatal("failed to send create data to proxy", err)
	}

	buf := make([]byte, PACKET_SIZE)
	n, err := proxy.Read(buf)
	if err != nil {
		log.Fatal("failed to read from proxy", err)
	}

	// log.Println(buf)
	if buf[0] != FLAG_SUCCESS {
		log.Fatal("failed to create minecraft server, please try again") // error messages to come
	}
	ip := string(buf[1:n])
	log.Println("Server now accessable please connect via", ip)

	go func() {
		_, err := io.Copy(dest, proxy)
		if err != nil {
			log.Println("connetion failed between proxy to dest", err)
		}
	}()

	_, err = io.Copy(proxy, dest)
	if err != nil {
		log.Println("connection failed between dest to proxy", err)
	}

	log.Println("server down | press any key to continue...")
	fmt.Scanf("")

	dest.Close()
	proxy.Close()
}

// package main
//
// import (
// 	"fmt"
// 	"io"
// 	"net"
// 	"time"
// )
//
// const DEST_IP = "192.168.1.145:25565"
// const HOST_IP = "127.0.0.1:25566"
//
// func main() {
// 	// Listen for incoming connections
// 	listener, err := net.Listen("tcp", HOST_IP)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer listener.Close()
//
// 	fmt.Println("Starting middleman loop")
// 	for {
// 		client, err := listener.Accept()
// 		if err != nil {
// 			fmt.Println("Error accepting client:", err)
// 			continue
// 		}
//
// 		// Create a new connection to the destination for each client
//
// 		go handleConn(client)
// 	}
// }
//
// func handleConn(client net.Conn) {
// 	fmt.Println("connection in coming")
// 	defer client.Close()
//
// 	dest, err := net.DialTimeout("tcp", DEST_IP, 30 * time.Second)
// 	if err != nil {
// 		fmt.Println("Error connecting to destination:", err)
// 		return
// 	}
// 	defer dest.Close()
//
// 	// capture connection ip
// 	buf := make([]byte, 1024)
// 	n, err := client.Read(buf)
// 	if err != nil {
// 		fmt.Println("failed to grab intial ip")
// 		return
// 	}
//
// 	// send back to dest so we don't drop packet and close connection
// 	dest.Write(buf[:n])
//
// 	// begin transport without anaylizing
//
// 	// Copy data from client -> dest
// 	go func() {
// 		_, err := io.Copy(dest, client)
// 		if err != nil {
// 			fmt.Println("Error copying from client to destination:", err)
// 		}
// 		// Closing dest when client disconnects
// 		dest.Close()
// 	}()
//
// 	// Copy data from dest -> client
// 	_, err = io.Copy(client, dest)
// 	if err != nil {
// 		fmt.Println("Error copying from destination to client:", err)
// 	}
// 	// Closing client when dest disconnects
// 	client.Close()
// }
