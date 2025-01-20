package packets

import (
	"encoding/binary"
)

func Create(flag uint8, id uint32, data []byte) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, id)

	// this operation will be kinda expensive
	// so might need to look into some time if it slows us down
	header := append([]byte{flag}, buf...)
	return append(header, data...)
}

// only call on KNOWN formated packets
// DOES NOT PROVIDE ANY WARNINGS AND WILL CRASH IF len(PACKET) < 5
// THIS SHOULD NEVER HAPPEN
func Read(n int, packet []byte) (uint8, uint32, []byte) {
	flag, data := packet[0], packet[5:n]
	id := binary.LittleEndian.Uint32(packet[1:5])

	return flag, id, data
}
