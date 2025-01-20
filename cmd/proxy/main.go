package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"sync"
	"sync/atomic"

	"github.com/minij147/proxycraft/pkg/consts"
	"github.com/minij147/proxycraft/pkg/packets"
)

var servers sync.Map

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
	"ear", "eat", "egg", "ego", "elf", "end", "era", "eve", "eye", "nnn",
}

type Server struct {
	Ip           string
	Conn         net.Conn
	NextClientID atomic.Uint32

	// not sync.map because each users is guranteed to have unique ID
	// given by the NextClientID.atomic, so threads (clients) will only
	// affect their given address
	Clients map[uint32]net.Conn
}

func ServerNew(conn net.Conn) *Server {
	return &Server{Conn: conn, Clients: make(map[uint32]net.Conn)}
}

// TODO: add checks for no duplicate ids
// WARNING: inital data max size should only be PACKET_SIZE_DATA
func (s *Server) AddClient(conn net.Conn, rawData []byte) bool {
	id := s.NextClientID.Add(1)
	_, ok := s.Clients[id]
	if ok {
		log.Println("id conflict on server")
		return false
	}
	s.Clients[id] = conn

	packet := packets.Create(consts.FLAG_CONNECTION_INCOMING, id, rawData)
	_, err := s.Conn.Write(packet)
	if err != nil {
		log.Println("failed write to server", err)
		return false
	}

	return true
}

// TODO: handle disconnections
func ClientRun(id uint32, client net.Conn, dest net.Conn) {
	// client --> server
	log.Println(id, "running")
	buf := make([]byte, consts.PACKET_SIZE_RAW)

	for {
		n, err := client.Read(buf)
		if err != nil {
			log.Println(id, "error from reading form client", err)
			break
		}

		packet := packets.Create(consts.FLAG_DATA, id, buf[:n])
		_, err = dest.Write(packet)
		if err != nil {
			log.Println(id, "error when writing to server", err)
			break
		}
	}

	log.Println(id, "ended connection")
	client.Close()
}

// Blocking, should be called in gorountinue
// starts servers runtime (anyalizes packets coming in and out)
func ServerRun(serv *Server) {
	log.Println("starting server runtime", serv.Ip)
	buf := make([]byte, consts.PACKET_SIZE_SIGNED)
	for {
		n, err := serv.Conn.Read(buf)
		if err != nil {
			log.Println("error encountered when reading packet", err, serv.Ip)
			break
		}

		flag, id, data := packets.Read(n, buf)
		// log.Println(flag, id, data)
		log.Println(buf[:n])
		log.Println(data)
		switch flag {
		case consts.FLAG_DATA:
			conn, ok := serv.Clients[id]
			if !ok { // idk what I want to do yet here...
				log.Println("invalid id")
				break
			}

			_, err := conn.Write(data)
			if err != nil {
				log.Println("invalid connection should remove")
			}
		case consts.FLAG_CONNECTION_ACCEPTED:
			log.Println(id, "connect has been accepted")
			conn, ok := serv.Clients[id]
			if !ok {
				packet := packets.Create(consts.FLAG_FAIL, id, []byte("failed to create client"))
				serv.Conn.Write(packet)
				conn.Close()
				break
			}

			go ClientRun(id, conn, serv.Conn)
		default:
			log.Println("not implemented flag")
		}
	}

	//TODO: add deletion and automactic disconnections for clients
	log.Println("stopping runtime for server")
}

func generateRandomIp(conn net.Conn) (string, bool) {
	serv := ServerNew(conn)
	for range consts.IP_GENERATE_CAP {
		i, j, k := rand.IntN(WORD_LEN), rand.IntN(WORD_LEN), rand.IntN(WORD_LEN)
		ip := fmt.Sprintf("%v.%v.%v.minics.dev", words[i], words[j], words[k])

		if _, loaded := servers.LoadOrStore(ip, serv); !loaded {
			log.Println(serv)
			return ip, true
		}
	}

	return "", false
}

/*
TODO:
1. Prevent duplicates of same connections
2. Add polling thread  which will poll connections to keep all alive
*/
func HandleServerCreation(flag uint8, client net.Conn) {
	fail := func() {
		client.Write([]byte{consts.FLAG_FAIL})
		client.Close()
	}

	clientIP := client.RemoteAddr().String()

	if flag != consts.FLAG_CREATE {
		log.Println("invalid byte code", clientIP, flag)

		client.Write([]byte("invalid code"))
		client.Close()
		return
	}

	log.Println("generating ip for", clientIP)
	ip, ok := generateRandomIp(client)
	if !ok {
		fail()
		return
	}

	servVal, ok := servers.Load(ip)
	if !ok {
		fail()
		return
	}

	serv, ok := servVal.(*Server)
	if !ok {
		// this should not happen
		log.Println("WARNING FAILED TRANSLATING TYPE for *Server")

		fail()
		return
	}

	go ServerRun(serv)

	msg := []byte("ip created: " + ip)
	packet := packets.Create(consts.FLAG_SUCCESS, 0, msg)
	client.Write(packet)

	log.Println("registered and now stored connection", clientIP, ip)
}

func HandleConnection(client net.Conn) {
	clientIP := client.RemoteAddr().String()

	// this is because it is are only unregistered packet
	// so we need to make room for header
	// also only header packet will contain no data (so no overflow)
	buf := make([]byte, consts.PACKET_SIZE_RAW)

	log.Println("rounting client", clientIP)
	n, err := client.Read(buf)
	if err != nil || n == 0 {
		log.Println("client failed to accept incoming backet", err)
		client.Close()
		return
	}

	log.Println(buf, n, clientIP)

	// creating server (loader connection packet)
	if n == consts.LOADER_CREATE_PACKET_SIZE {
		HandleServerCreation(buf[0], client)
		return
	}

	// creating user (user connection packet)
	if n < consts.IP_SIZE {
		log.Println("client invalid size for ip")
		client.Close()
		return
	}

	requestIP := string(buf[5:consts.IP_SIZE])

	destVal, ok := servers.Load(requestIP)
	if !ok {
		log.Println("requested ip does not exists", requestIP, clientIP)
		client.Close()
		return
	}

	dest, ok := destVal.(*Server)
	if !ok {
		//warning because this should not fail ever
		log.Println("WARNING: failed to cast connection, SERVER SIDE ERROR")
		client.Close()
		return
	}
	log.Println(dest)

	if !dest.AddClient(client, buf[:n]) {
		log.Println("FAILED TO ADD CLIENT TO SERVER")
		client.Close()
	}

	// start the initalization of the client
	// aka: creating on the loaders side a valid connection or returning an already exisiting on
	// this is done through the user ids

	// TODO: send ip along with it so the loader can check connections and ips and stuff
	// _, err = dest.Conn.Write([]byte{consts.FLAG_CONNECTION_NEW, 0})
	// if err != nil {
	// 	// TODO: depending on error we will have to close the connection
	// 	log.Println("failed to write connection new packet", err)
	// 	client.Close()
	// 	return
	// }

	// dest won't directly write a response but
	/*
		go func(){
			packet = {FLAG, ID, DATA}

			packet = server.Read()
			id := packet[1]
			client[id].Write(packet)}
	*/

}

func main() {
	log.Println("starting the proxy")

	listener, err := net.Listen("tcp", consts.IP_PROXY)
	if err != nil {
		log.Fatal("failed to start listening", err, consts.IP_PROXY)
	}

	for {
		client, err := listener.Accept()
		if err != nil {
			log.Println("failed to accept incoming connection", err)
			continue
		}

		go HandleConnection(client)
	}
}
