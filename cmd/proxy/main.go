package main

import (
	"log"
	"os"
	"time"

	"github.com/minij147/proxycraft/internal/proxy/services"
	"github.com/minij147/proxycraft/pkg/consts"
)

/*
Proxy.go

Proxy.go is to run the proxy which allows for loaders and clients to connect and communicate with port forwarding.

Runs the proxy server were the main logic is contained here and runs a http server for config details for clients.

When running proxy you must provide the following arguments (in the order as appeared)

{connection ip} - This will be the domain name or whatever you set for the loaders to connect too
{connection port} - This will be the port used to reach the proxy from the outside (xxx.minics.dev:{connection-port})
{proxy host port} - This will be the port the tcp server will listen on
{config host port} - This will be the port for the http config server will listen on
*/

const ARGS_SIZE = 4

func main() {
	log.Println(os.Args)
	if len(os.Args) < ARGS_SIZE+1 {
		log.Fatal("not enough arguments passed in, please provide the following:\n{connection-ip} {connection port} {proxy host port} {config host port}\nSee docs for more info")
	}

	connectionIP := os.Args[1]
	connectionPort := os.Args[2]
	proxyIp := ":" + os.Args[3]
	configIp := ":" + os.Args[4]

	cfg := consts.ServerConfig{
		ProxyIP:             connectionIP,
		ProxyPort:           connectionPort,
		ServerVersion:       consts.VERSION_SERVER,
		LatestClientVersion: consts.VERSION_CLIENT,
		LastDeployedDate:    time.Now().Format("2006-01-02 15:04:05"),
	}

	log.Println("Starting Proxy on Version", consts.VERSION_SERVER)
	log.Printf("\nINFO:\nConnection Ip: %v\nConnection Port: %v\nProxy Host: %v\nConfig Host: %v\n\n", connectionIP, connectionPort, proxyIp, configIp)
	services.ConfigRun(configIp, cfg)

	services.ProxyRun(proxyIp)
}
