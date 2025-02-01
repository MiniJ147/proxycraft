package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/minij147/proxycraft/pkg/consts"
)

//TODO: add debug mode

const PROXY_URL = "https://proxycraft.minics.dev/config"
// const PROXY_URL = "http://localhost:3002/config"

var IP_DEST = ":25565"
var IP_PROXY_CONN = ""
var URL = ""

func NewClient(prox net.Conn, ip string) {
	dest, e := net.Dial("tcp", IP_DEST)
	if e != nil {
		log.Println("failed to connect to dest", e)
		return
	}
	defer dest.Close()

	log.Println(ip)
	pipe, e := net.Dial("tcp", IP_PROXY_CONN)
	if e != nil {
		log.Println("failed to create pipe", e)
		return
	}
	defer pipe.Close()

	msg := []byte(URL)
	_, e = pipe.Write(append([]byte{consts.FLAG_CONN_OK}, msg...))
	if e != nil {
		log.Println("failed to write throuhgh pipe")
		return
	}

	go func() {
		_, e := io.Copy(dest, pipe)
		if e != nil {
			log.Println(e)
		}
	}()

	_, e = io.Copy(pipe, dest)
	if e != nil {
		log.Println(e)
	}
}

func FetchConfig() consts.ServerConfig {
	var cfg consts.ServerConfig
	res, err := http.Get(PROXY_URL)
	if err != nil {
		log.Fatalf("failed to config from %v with error %v\nplease report for more information", PROXY_URL, err)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("failed to read body data config from %v with error %v\nplease report for more information", PROXY_URL, err)
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("failed to unmarshal data from config with error %v\nplease report for more information", err)
	}

	return cfg
}

func getPort() string {
	const DEFAULT = "25565"

	var port string
	fmt.Print("Enter your port for your server (default 25565), enter nothing for default: ")
	_, e := fmt.Scanln(&port)
	if e != nil {
		log.Println("no input defaulting port")
		return DEFAULT
	}

	_, e = strconv.Atoi(port)
	if e != nil {
		log.Printf("input failed check with error %v\ndefaulting port")
		return DEFAULT
	}

	return port
}

func main() {
	log.Println("starting loader...")
	cfg := FetchConfig()

	IP_PROXY_CONN = fmt.Sprintf("%v:%v", cfg.ProxyIP, cfg.ProxyPort)
	IP_DEST = ":" + getPort()

	log.Printf("Hosting on 127.0.0.1:%v | Connecting to proxy: %v\n", IP_DEST, IP_PROXY_CONN)

	p, e := net.Dial("tcp", IP_PROXY_CONN)
	if e != nil {
		log.Println("failed to connection", e)
	}

	_, e = p.Write([]byte{consts.FLAG_INIT})
	if e != nil {
		log.Println("failed to write", e)
	}

	buf := make([]byte, consts.PACKET_SIZE)

	n, e := p.Read(buf)
	if e != nil {
		log.Fatal("failed to get response")
	}

	if buf[0] != consts.FLAG_INIT_OK {
		log.Fatal("failed to init")
	}

	URL = string(buf[1:n])
	log.Println("initlized with server", URL)
	log.Println(p.RemoteAddr().String())

	for {
		n, e := p.Read(buf)
		if e != nil {
			log.Fatal("failed to read to buf", e)
		}

		switch buf[0] {
		case consts.FLAG_CONN_NEW:
			go NewClient(p, string(buf[1:n]))
		case consts.FLAG_POLL:
		default:
			log.Println("unkown flag sent")
		}
	}
}
