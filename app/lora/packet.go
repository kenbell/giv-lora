package lora

import (
	"encoding/binary"
	"fmt"

	"github.com/kenbell/giv-lora/app/modbus"
)

// Packets look something like this:
//    88 41   3b    c5 1a ff ff 10 01    66    b8 bb
//    sig     cnt   fixed                data  chksum
//
// sig    - fixed signature
// cnt    - incrementing counter (per packet)
// fixed  - static data, purpose unknown
// data   - the actual payload
// chksum - MODBUS-16 CRC of the entire packet (including sig)

// ParseLoRaPacket validates and extracts fields from the packet
//
// On success, the packet sequence number and payload are returned.
// Otherwise an error indicates the packet failed validation.
func ParseLoRaPacket(pkt []byte) (int, []byte, error) {
	if len(pkt) < 11 {
		return 0, nil, fmt.Errorf("invalid packet (too small)")
	}

	if pkt[0] != 0x88 || pkt[1] != 0x41 {
		return 0, nil, fmt.Errorf("invalid packet (bad sig)")
	}

	calcCrc := modbus.CRC(pkt[:len(pkt)-2])
	pktCrc := binary.BigEndian.Uint16(pkt[len(pkt)-2:])

	if calcCrc != pktCrc {
		return 0, nil, fmt.Errorf("invalid packet (bad crc)")
	}

	cnt := int(pkt[2])
	payload := pkt[9 : len(pkt)-2]

	return cnt, payload, nil
}
