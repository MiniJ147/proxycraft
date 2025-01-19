package main

// first data sent by minecraft is connection address

import (
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net"
	"sync"
	"time"
)

const WORD_LEN = 100
const IP_SIZE = 27 // aaa.bbb.ccc.minics.dev:XXXX
const MIDDLE_IP = "127.0.0.1:25565"
const CREATE_BYTES uint8 = 201

// required becuase minecraft sends 5 extra bytes which I don't know that they are
// so we must skip past that to read the domain name
const INITAL_PACKET_OFFSET = 5
const PACKET_SIZE = 1024
const TIMEOUT_INTIAL = 30 * time.Second
const TIMEOUT_WRITE = 30 * time.Second

const FLAG_CREATE = 100
const FLAG_SUCCESS = 101
const FLAG_FAIL = 102

const RETRY_IP_GENERATE_CAP = 20

var servers sync.Map

type Server struct {
	IP string
}

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
	"ear", "eat", "egg", "ego", "elf", "end", "era", "eve", "eye", "nnn",
}

func generateRandomIp(hostIp string) (string, bool) {
	for range RETRY_IP_GENERATE_CAP {
		i, j, k := rand.IntN(WORD_LEN), rand.IntN(WORD_LEN), rand.IntN(WORD_LEN)
		ip := fmt.Sprintf("%v.%v.%v.minics.dev", words[i], words[j], words[k])

		if _, loaded := servers.LoadOrStore(ip, hostIp); !loaded {
			return ip, true
		}
	}

	return "", false
}

func BeginExchange(client net.Conn, dest net.Conn, serverGenerateIp string) {
	defer client.Close()
	defer dest.Close()

	clientIp := client.RemoteAddr().String()
	destIp := client.RemoteAddr().String()

	log.Println("beginning to exchange between", clientIp, destIp, serverGenerateIp)

	go func() {
		_, err := io.Copy(dest, client)
		if err != nil {
			log.Println("connetion failed between client to dest", err)
		}
	}()

	_, err := io.Copy(client, dest)
	if err != nil {
		log.Println("connection failed between dest to client", err)
	}
}

func RouteClient(client net.Conn) {
	clientIp := client.RemoteAddr().String()
	buf := make([]byte, PACKET_SIZE)

	log.Println("routing client", clientIp)
	n, err := client.Read(buf)
	if err != nil || n == 0 {
		log.Println("client failed to accept incoming packet", err)
		client.Close()
		return
	}

	// log.Println(buf, n, clientIp)

	// if a client is attempting to connect to a server
	// not create one
	if n >= IP_SIZE {
		requestedIP := string(buf[5:IP_SIZE])
		log.Println("connection for ", requestedIP, "from", clientIp)

		ipVal, ok := servers.Load(requestedIP)
		if !ok {
			log.Println("request ip does not exists", requestedIP, clientIp)
			client.Close()
			return
		}

		ip, ok := ipVal.(string)
		if !ok {
			log.Println("WARNING FAILED TO TRANSLATE IP, INTENRAL ISSUE")
			client.Close()
			return
		}

		log.Println(clientIp, "found", ip, requestedIP)

		dest, err := net.DialTimeout("tcp", ip, TIMEOUT_INTIAL)
		if err != nil {
			log.Println("connection for ", requestedIP, "from", clientIp, "failed")
			client.Close()
			return
		}

		_, err = dest.Write(buf[:n])
		if err != nil {
			log.Println("failed to write packet back to dest", requestedIP, clientIp)
			dest.Close()
			client.Close()
			return
		}

		go BeginExchange(client, dest, requestedIP)
		return
	}

	// invalid case
	if n != 1 {
		log.Println("invalid packet size, below min for ip and not one byte for creation flag", clientIp, n)
		client.Close()
		return
	}

	if buf[0] != FLAG_CREATE {
		log.Println("invalid byte code", clientIp, buf[0])
		client.Write([]byte("invalid code"))
		client.Close()
		return
	}

	log.Println("creating server for", clientIp)
	ip, success := generateRandomIp(clientIp + "25566")
	if !success {
		client.Write([]byte{FLAG_FAIL})
	} else {
		msg := []byte("ip created: " + ip)
		client.Write(append([]byte{FLAG_SUCCESS}, msg...))
	}

	client.Close()
}

func main() {
	log.Println("hello from middleman")
	middleman, err := net.Listen("tcp", MIDDLE_IP)
	if err != nil {
		panic(err)
	}
	defer middleman.Close()

	for {
		client, err := middleman.Accept()
		if err != nil {
			log.Println("client failed to connect", err)
			continue
		}

		go RouteClient(client)
	}
}
