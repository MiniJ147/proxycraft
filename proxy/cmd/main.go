package main

// first data sent by minecraft is connection address

import (
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net"
	"strings"
	"sync"
)

var ports sync.Map
var servers sync.Map

const PORT_LOW = 49152
const PORT_HIGH = 65535
const PORT_RANGE = PORT_HIGH - PORT_LOW
const WORD_LEN = 100
const IP_SIZE = 27 // aaa.bbb.ccc.minics.dev:XXXX
const PACKET_SIZE = 1024

// required becuase minecraft sends 5 extra bytes which I don't know that they are
// so we must skip past that to read the domain name
const MINECRAFT_DOMAIN_OFFSET = 5

var words = []string{
	"ace", "act", "add", "ado", "aft", "age", "ago", "aid", "air", "ale",
	"all", "and", "ant", "any", "ape", "apt", "arc", "arm", "art", "ash",
	"ask", "ate", "aub", "awe", "axe", "bad", "bag", "ban", "bar", "bat",
	"bay", "bed", "bee", "beg", "bet", "bib", "big", "bin", "bit", "bog",
	"boo", "bow", "box", "boy", "bud", "bug", "bun", "bus", "but", "buy",
	"cab", "cad", "cam", "can", "cap", "car", "cat", "cay", "cod", "cog",
	"con", "coo", "cop", "cot", "cow", "coy", "cry", "cub", "cue", "cup",
	"cur", "cut", "dab", "dad", "dam", "day", "den", "dew", "did", "die",
	"dig", "dim", "din", "dip", "dot", "dry", "dub", "dug", "due", "dye",
	"ear", "eat", "egg", "ego", "elf", "end", "era", "eve", "eye",
}

type Server struct {
	Host net.Conn
	Port int
}

func InitializerListen() {
	tcp, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatalf("Failed creating Initializer Listener on port 3000 | %v", err)
	}
	defer tcp.Close()

	log.Println("listening on port 3000")

	for {
		conn, err := tcp.Accept()
		if err != nil {
			log.Println("failed to accept incoming connection", err)
			continue
		}

		go InitializerHandleConnection(conn)
	}
}

func InitializerHandleConnection(conn net.Conn) {
	for {
		newIP := generateRandomIp()
		if _, loaded := servers.LoadOrStore(newIP, conn); !loaded {
			fmt.Println("created ip", newIP)
			conn.Write([]byte(newIP))
			break
		}
	}
}

func generateRandomIp() string {
	i, j, k := rand.IntN(WORD_LEN), rand.IntN(WORD_LEN), rand.IntN(WORD_LEN)
	return fmt.Sprintf("%v.%v.%v.minics.dev", words[i], words[j], words[k])
}

func TunnelHandleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, PACKET_SIZE)
	size, err := conn.Read(buf)
	if err != nil || size == 0 {
		log.Println(conn.RemoteAddr().String(), "conn died")
	}

	connIP := strings.TrimSpace(string(buf[MINECRAFT_DOMAIN_OFFSET:IP_SIZE]))
	fmt.Println(connIP, []byte(connIP))

	val, ok := servers.Load(connIP)
	if !ok {
		log.Println("cannot find ip")
		return
	}

	host, ok := val.(net.Conn)
	if !ok {
		log.Println("failed to translate type out of server map")
		return
	}

	// send the packet we caught and anaylzed to the server
	host.Write(buf[:size])
	go TunnelHostBack(connIP, conn, host)
    
    fmt.Println("starting to listen for our packets")
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("connection done")
			break
		}

		_, err = host.Write(buf[:n])
		if err != nil {
			if err == io.EOF {
				servers.Delete(connIP)
				log.Println("host done")
			}

			log.Println("host closed")
			break
		}
	}
}

func TunnelHostBack(ip string, conn net.Conn, host net.Conn) {
	buf := make([]byte, PACKET_SIZE)

	for {
		n, err := host.Read(buf)
		if err != nil {
			if err == io.EOF || n == 0 {
				servers.Delete(ip)
				log.Println("host done")
			}
			log.Println("connection done")
			break
		}

		_, err = conn.Write(buf[:n])
		if err != nil {
			break
		}
	}
}

func TunnelListen() {
	tcp, err := net.Listen("tcp", "127.0.0.1:3001")
	if err != nil {
		log.Fatalf("failed to create tunneler on 3001 | %v", err)
	}
	defer tcp.Close()

	for {
		conn, err := tcp.Accept()
		if err != nil {
			log.Println("failed to accept incoming conncetion to middle man", err)
			continue
		}

		go TunnelHandleConnection(conn)
	}
}

func main() {
	fmt.Println("Hello World!")
	go InitializerListen()
	TunnelListen()
}
