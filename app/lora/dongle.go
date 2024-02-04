package lora

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"go.bug.st/serial"
)

type Dongle struct {
	file io.ReadWriteCloser
}

func NewDongle(path string) (*Dongle, error) {
	port, err := serial.Open(path, &serial.Mode{
		BaudRate: 115200,
	})

	if err != nil {
		return nil, err
	}

	// reset the dongle state
	_, err = port.Write([]byte("\nreset\n"))
	if err != nil {
		return nil, err
	}

	time.Sleep(10 * time.Millisecond)
	port.ResetInputBuffer()

	return &Dongle{file: port}, nil
}

func (d *Dongle) StartSniff(channel int) error {
	_, err := d.file.Write([]byte(fmt.Sprintf("sniff %d\n", channel)))
	return err
}

func (d *Dongle) Read(buf []byte) (int, error) {
	hdr := []byte{0, 0}
	_, err := io.ReadFull(d.file, hdr)
	if err != nil {
		return 0, err
	}

	n := binary.BigEndian.Uint16(hdr)

	return io.ReadFull(d.file, buf[:n])
}
