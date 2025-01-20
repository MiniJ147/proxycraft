package consts

const (
	IP_PROXY        = "127.0.0.1:25565"
	IP_SIZE         = 27
	IP_GENERATE_CAP = 20

	// size of standard data
	PACKET_SIZE_DATA = 1024

	// size for our id
	PACKET_SIZE_ID = 4 

	//size for our flag
	PACKET_SIZE_FLAG = 1

	// total size with layout
	PACKET_SIZE = PACKET_SIZE_FLAG + PACKET_SIZE_ID + PACKET_SIZE_DATA

	// these  are 1 bytes flags 0-255
	FLAG_CREATE                uint8 = 100
	FLAG_SUCCESS               uint8 = 101
	FLAG_FAIL                  uint8 = 102
	FLAG_POLL                  uint8 = 103 // polls connection (disregard for loader)
	FLAG_CONNECTION_DISCONNECT uint8 = 104
	FLAG_DATA                  uint8 = 105
	FLAG_CONNECTION_NEW        uint8 = 106

	LOADER_CREATE_PACKET_SIZE = 1
)

// DO NOT MODIFY
var PACKET_EMPTY = []byte{}
