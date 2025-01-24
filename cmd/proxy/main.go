package main

import (
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/minij147/proxycraft/pkg/consts"
)

const WORD_LEN = 100

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
	"ear", "eat", "egg", "ego", "elf", "aid", "ann", "lol", "anf", "ade",
}

var m map[string]net.Conn = make(map[string]net.Conn)
var pipes []net.Conn = make([]net.Conn, 0)

func HandleLoaderInit(conn net.Conn, ip string) {
	_, ok := m[ip]
	if ok {
		log.Println("cannot create server already exists")
		conn.Write([]byte{consts.FLAG_INIT_FAIL})
		return
	}

	m[ip] = conn
	m[consts.TEST_IP] = conn

	msg := []byte(consts.TEST_IP)
	conn.Write(append([]byte{consts.FLAG_INIT_OK}, msg...))
}

func HandleClientJoin(conn net.Conn, ip string, payload []byte, n int) {
	if n < consts.IP_SIZE {
		conn.Close()
		return
	}

	url := string(payload[5:consts.IP_SIZE])
	log.Println(ip, "-->", url)

	serv, ok := m[url]
	if !ok {
		log.Println("server does not exists")
		conn.Close()
		return
	}

	ipBytes := []byte(ip)
	_, e := serv.Write(append([]byte{consts.FLAG_CONN_NEW}, ipBytes...))
	if e != nil {
		log.Fatal("failed to write to server")
		conn.Close()
		return
	}

	// spin until pipe created
	log.Println("waiting...")
	for len(pipes) == 0 {
		time.Sleep(500 * time.Millisecond)
	}
	log.Println("got pipe...")

	pipe := pipes[0]
	pipes = pipes[1:]

	_, e = pipe.Write(payload)
	if e != nil {
		log.Println("failed wirte", e)
	}

	go func() {
		_, e := io.Copy(conn, pipe)
		if e != nil {
			log.Println("pipe -> conn", e)
		}
	}()

	_, e = io.Copy(pipe, conn)
	if e != nil {
		log.Println("conn -> pipe", e)
	}
}

func HandleConnection(conn net.Conn) {
	ip := strings.Split(conn.RemoteAddr().String(), ":")[0]

	buf := make([]byte, consts.PACKET_SIZE)
	n, e := conn.Read(buf)
	if e != nil {
		conn.Close()
		log.Println("failed to read client buf", e)
		return
	}

	if n > 1 {
		// client trying to join minecraft server
		log.Println("joining minecraft server")
		HandleClientJoin(conn, ip, buf, n)
		return
	}

	switch buf[0] {
	case consts.FLAG_INIT:
		log.Println("initlizing server")
		HandleLoaderInit(conn, ip)
	case consts.FLAG_CONN_OK:
		log.Println("found connection")
		pipes = append(pipes, conn)
	default:
		log.Println("invalid switch")
		break
	}

}

func main() {
	log.Println("starting proxy...")

	l, e := net.Listen("tcp", consts.IP_PROXY)
	if e != nil {
		log.Fatal("failed to start server on ", consts.IP_PROXY)
	}

	for {
		c, e := l.Accept()
		if e != nil {
			log.Fatal(e)
		}
		log.Println(m)

		go HandleConnection(c)
	}
}
