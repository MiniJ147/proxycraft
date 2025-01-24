package consts

const (
	IP_PROXY_HOST = ":25565"
	IP_PROXY_CONN   = "proxy.minics.dev:26850"
	// IP_PROXY_CONN   = "127.0.0.1:25565"
	IP_SIZE         = 27
	IP_GENERATE_CAP = 20

	PACKET_SIZE = 1024

	TEST_IP = "aaa.bbb.ccc.minics.dev"

	FLAG_INIT      uint8 = 1
	FLAG_INIT_OK   uint8 = 2
	FLAG_INIT_FAIL uint8 = 3
	FLAG_CONN_NEW  uint8 = 10
	FLAG_CONN_OK   uint8 = 11
	FLAG_CONN_FAIL uint8 = 12
	FLAG_POLL      uint8 = 20
)
