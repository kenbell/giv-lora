package modbus

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type ADU struct {
	Address  byte
	Function byte
	Payload  []byte
}

func ParseADU(buf []byte) (int, *ADU, error) {
	if len(buf) < 4 {
		return 0, nil, fmt.Errorf("buffer too small")
	}

	adu := &ADU{
		Address:  buf[0],
		Function: buf[1],
	}

	len := packetLength(buf)
	if len < 4 {
		return 0, nil, fmt.Errorf("no valid modbus ADU detected")
	}

	// Strip addr/fn from start and checksum from end
	adu.Payload = make([]byte, len-4)
	copy(adu.Payload, buf[2:len-2])

	return len, adu, nil
}

func (adu *ADU) Dump() string {
	return fmt.Sprintf("Addr: %d, Fn: %d, Payload: %s", adu.Address, adu.Function, hex.EncodeToString(adu.Payload))
}

func packetLength(buf []byte) int {
	// We don't know if a packet is a request or a response, so we scan
	// the buffer looking for a valid checksum (indicating end of packet)
	for l := 2; l <= len(buf)-2; l++ {
		calcCrc := CRC(buf[:l])
		if binary.LittleEndian.Uint16(buf[l:]) == calcCrc {
			return l + 2
		}
	}

	return -1
}
