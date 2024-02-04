package modbus

import (
	"encoding/binary"
	"fmt"
)

type Device struct {
	Address byte

	InputRegisters map[uint16]uint16

	lastADU *ADU
}

func NewDevice(address byte) *Device {
	return &Device{
		Address:        address,
		InputRegisters: map[uint16]uint16{},
	}
}

func (d *Device) SetInputRegister(reg uint16, val uint16) {
	d.InputRegisters[reg] = val
}

// HandleADU updates the device state based on observed ADUs.
//
// ADUs are request/response and a valid pair is required to correctly
// interpret the modbus transaction, so requests are stored temporarily.
//
// The number of
func (d *Device) HandleADU(adu *ADU) error {
	if adu.Address != d.Address {
		return fmt.Errorf("incorrect device")
	}

	// Currently, only input registers are modelled
	if adu.Function != 4 {
		d.lastADU = nil
		return fmt.Errorf("unhandled function: %d", adu.Function)
	}

	if len(adu.Payload) == 4 {
		// infer request, stash it
		d.lastADU = adu
		return nil
	}

	if len(adu.Payload) < 1 || len(adu.Payload) != int(adu.Payload[0])+1 {
		// malformed since first byte should be payload length in response
		d.lastADU = nil
		return fmt.Errorf("malformed modbus response packet")
	}

	if d.lastADU == nil {
		return fmt.Errorf("response ADU with for unknown request")
	}

	//
	// To get here: we have a response and it matches an immediately preceeding
	// request.
	//

	start := binary.BigEndian.Uint16(d.lastADU.Payload[0:])
	count := binary.BigEndian.Uint16(d.lastADU.Payload[2:])

	if uint16(adu.Payload[0]) != count*2 {
		d.lastADU = nil
		return fmt.Errorf("response ADU has unexpected size")
	}

	for reg := uint16(0); reg < count; reg++ {
		//fmt.Printf("pl: %s", hex.EncodeToString(adu.Payload))
		v := binary.BigEndian.Uint16(adu.Payload[1+2*reg:])
		d.InputRegisters[start+reg] = v
	}

	return nil
}
