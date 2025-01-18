// package main
//
// import (
// 	"fmt"
// 	"net"
// )
//
// const DEST_IP = "192.168.1.145:25565"
// const HOST_IP = "127.0.0.1:3000"
//
// func main() {
// 	mid, err := net.Listen("tcp", HOST_IP)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	dest, err := net.Dial("tcp", DEST_IP)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	fmt.Println("starting middle man loop")
// 	for {
// 		conn, err := mid.Accept()
// 		if err != nil {
// 			fmt.Println("Error:", err)
// 			continue
// 		}
//
// 		go handleConn(conn, dest)
// 	}
// }
//
// func handleConn(client net.Conn, dest net.Conn) {
// 	defer client.Close()
//
// 	// handles dest --> client
// 	go func() {
// 		buf := make([]byte, 1024)
//
// 		for {
// 			n, err := dest.Read(buf)
// 			if err != nil {
// 				fmt.Println("error from dest read", err)
// 				break
// 			}
//
// 			_, err = client.Write(buf[:n])
// 			if err != nil {
// 				fmt.Println("write to clinet failed", err)
// 				break
// 			}
// 		}
//
// 		fmt.Println("killed packets sending to client")
// 	}()
//
// 	// handles client -> dest
// 	buf := make([]byte, 1024)
// 	for {
// 		n, err := client.Read(buf)
// 		if err != nil {
// 			fmt.Println("error from clinet read", err)
// 			break
// 		}
//
// 		_, err = dest.Write(buf[:n])
// 		if err != nil {
// 			fmt.Println("write from client to dest failed", err)
// 			break
// 		}
// 	}
//
// 	fmt.Println("killed packets sending to dest")
// }
package main

import (
	"fmt"
	"io"
	"net"
)

const DEST_IP = "192.168.1.145:25565"
const HOST_IP = "127.0.0.1:3000"

func main() {
	// Listen for incoming connections
	listener, err := net.Listen("tcp", HOST_IP)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Starting middleman loop")
	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting client:", err)
			continue
		}

		// Create a new connection to the destination for each client
		dest, err := net.Dial("tcp", DEST_IP)
		if err != nil {
			fmt.Println("Error connecting to destination:", err)
			client.Close()
			continue
		}

		go handleConn(client, dest)
	}
}

func handleConn(client net.Conn, dest net.Conn) {
	defer client.Close()
	defer dest.Close()

	// Copy data from client -> dest
	go func() {
		_, err := io.Copy(dest, client)
		if err != nil {
			fmt.Println("Error copying from client to destination:", err)
		}
		// Closing dest when client disconnects
		dest.Close()
	}()

	// Copy data from dest -> client
	_, err := io.Copy(client, dest)
	if err != nil {
		fmt.Println("Error copying from destination to client:", err)
	}
	// Closing client when dest disconnects
	client.Close()
}

