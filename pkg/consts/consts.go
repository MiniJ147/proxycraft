package consts

const (
	IP_PROXY        = "127.0.0.1:25565"
	IP_SIZE         = 27
	IP_GENERATE_CAP = 20

	// size of standard data
	// size of responses from servers (exlcudes header information)
	PACKET_SIZE_RAW = 2048

	// size for our id
	PACKET_SIZE_ID = 4

	//size for our flag
	PACKET_SIZE_FLAG = 1
	PACKET_SIZE_HEADER_OFFSET = PACKET_SIZE_FLAG + PACKET_SIZE_ID

	// total size with layout
	PACKET_SIZE_SIGNED = PACKET_SIZE_HEADER_OFFSET+ PACKET_SIZE_RAW

	// these  are 1 bytes flags 0-255
	FLAG_CREATE                uint8 = 100
	FLAG_SUCCESS               uint8 = 101
	FLAG_FAIL                  uint8 = 102
	FLAG_POLL                  uint8 = 103 // polls connection (disregard for loader)
	FLAG_DATA                  uint8 = 104
	FLAG_CONNECTION_DISCONNECT uint8 = 105
	FLAG_CONNECTION_INCOMING   uint8 = 106
	FLAG_CONNECTION_ACCEPTED   uint8 = 107
	FLAG_CONNECTION_FAILED     uint8 = 108

	LOADER_CREATE_PACKET_SIZE = 1
)

// DO NOT MODIFY
var PACKET_EMPTY = []byte{}
