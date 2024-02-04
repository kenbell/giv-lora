package modbus

func CRC(buf []byte) uint16 {
	crc := uint16(0xFFFF)

	for _, b := range buf {
		crc ^= uint16(b)

		for bit := 0; bit < 8; bit++ {
			if (crc & 0x0001) != 0 {
				crc >>= 1
				crc ^= 0xA001
			} else {
				crc >>= 1
			}
		}
	}

	return crc
}
