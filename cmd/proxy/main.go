package main

import (
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net"
	"strings"
	"sync"
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

// var m map[string]net.Conn = make(map[string]net.Conn)

// var pipes []net.Conn = make([]net.Conn, 0)
var servers sync.Map

type Server struct {
	Conn     net.Conn
	Ip       string
	IpCustom string
	PipeChan chan (net.Conn)
}

func ServerNew(conn net.Conn, ip string) *Server {
	return &Server{
		Conn:     conn,
		Ip:       ip,
		PipeChan: make(chan net.Conn),
	}
}

func LoadIntoMap(serv *Server) (string, bool) {
	_, loaded := servers.LoadOrStore(serv.Ip, serv)
	if loaded {
		return serv.IpCustom, false
	}

	for range consts.IP_GENERATE_CAP {
		i, j, k := rand.IntN(WORD_LEN), rand.IntN(WORD_LEN), rand.IntN(WORD_LEN)
		ip := fmt.Sprintf("%v.%v.%v.minics.dev", words[i], words[j], words[k])

		if _, loaded := servers.LoadOrStore(ip, serv); !loaded {
			log.Println(serv, "added")
			serv.IpCustom = ip
			return ip, true
		}
	}

	return "", false
}

func RemoveFromMap(serv *Server) {
	servers.Delete(serv.Ip)
	servers.Delete(serv.IpCustom)
}

func HandleLoaderInit(conn net.Conn, ip string) {
	_, ok := servers.Load(ip)
	if ok {
		log.Println("cannot create server already exists")
		conn.Write([]byte{consts.FLAG_INIT_FAIL})
		return
	}

	serv := ServerNew(conn, ip)

	log.Println("hosting on", ip)

	// servers.Store(ip, serv)
	// servers.Store(consts.TEST_IP, serv)
	ipCustom, ok := LoadIntoMap(serv)
	if !ok {
		log.Println("failed to write into map", ipCustom)
		conn.Write([]byte{consts.FLAG_INIT_FAIL})
		return
	}

	msg := []byte(ipCustom)
	conn.Write(append([]byte{consts.FLAG_INIT_OK}, msg...))

	// polling functions
	go func() {
		for {
			_, err := conn.Write([]byte{consts.FLAG_POLL})
			if err != nil {
				log.Println("failed polling should kill connection")
				break
			}
			time.Sleep(5 * time.Second)
		}

		log.Println("removed", serv.Ip, serv.IpCustom)
		RemoveFromMap(serv)
	}()

}

func HandleClientJoin(conn net.Conn, ip string, payload []byte, n int) {
	if n < consts.IP_SIZE {
		conn.Close()
		return
	}

	url := string(payload[5:consts.IP_SIZE])

	if !strings.Contains(url, ".minics.dev") {
		url = "127.0.0.1"
	}
	log.Println(ip, "-->", url)

	servVal, ok := servers.Load(url)
	if !ok {
		log.Println("server does not exists")
		conn.Close()
		return

	}

	serv, ok := servVal.(*Server)
	if !ok {
		log.Println("WARNING FAILED TYPECAST THIS SHOULD NOT HAPPEN")
		conn.Close()
		return
	}

	ipBytes := []byte(ip)
	_, e := serv.Conn.Write(append([]byte{consts.FLAG_CONN_NEW}, ipBytes...))
	if e != nil {
		log.Fatal("failed to write to server")
		conn.Close()
		return
	}

	// spin until pipe created
	log.Println("waiting...")
	pipe := <-serv.PipeChan
	log.Println("got pipe...")

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

	pipe.Close()
	conn.Close()
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
		servVal, ok := servers.Load(ip)
		if !ok {
			log.Println("failed to find server")
			return
		}

		serv, ok := servVal.(*Server)
		if !ok {
			log.Println("FAILED TYPECASE IN FLAG CONN OK")
			return
		}

		serv.PipeChan <- conn
	default:
		log.Println("invalid switch")
	}

}

func main() {
	log.Println("starting proxy...")

	l, e := net.Listen("tcp", consts.IP_PROXY_HOST)
	if e != nil {
		log.Fatal("failed to start server on ", consts.IP_PROXY_HOST)
	}

	for {
		c, e := l.Accept()
		if e != nil {
			log.Fatal(e)
		}
		log.Println(servers)

		go HandleConnection(c)
	}
}
