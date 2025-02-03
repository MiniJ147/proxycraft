package consts

import "time"

type ServerConfig struct {
	ProxyIP             string
	ProxyPort           string
	ServerVersion       string
	LatestClientVersion string
	LastDeployedDate    string
}

const (
	VERSION_SERVER = "v0.1"
	VERSION_CLIENT = "v0.1"

	// IP_PROXY_HOST        = ":25565"
	// IP_PROXY_CONFIG_HOST = ":3000"
	// IP_PROXY_CONN        = "proxy.minics.dev:26850"
	// IP_PROXY_CONN   = "127.0.0.1:25565"

	// size of ip for our minecraft packet
	IP_CLIENT_SIZE  = 27
	IP_GENERATE_CAP = 20

	// length of generated ips len(xxx.xxx.xxx.minics.dev)
	IP_LEN = 22

	PACKET_SIZE = 1024

	FLAG_INIT      uint8 = 1
	FLAG_INIT_OK   uint8 = 2
	FLAG_INIT_FAIL uint8 = 3
	FLAG_CONN_NEW  uint8 = 10
	FLAG_CONN_OK   uint8 = 11
	FLAG_CONN_FAIL uint8 = 12
	FLAG_POLL      uint8 = 20

	POLL_FREQUENCY = 30 * time.Second

	DEBUG_DOMAIN = "deb.ugx.xxx.minics.dev"
)
