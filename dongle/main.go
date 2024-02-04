package main

import (
	"machine"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kenbell/tinygo-lora/drivers/sx127x"
)

var (
	spi = machine.SPI0
	d0  = machine.GPIO21
	rst = machine.GPIO20
	cs  = machine.GPIO17
)

var (
	radio *sx127x.Device
)

var (
	radioWaitGroup sync.WaitGroup
	stopRadio      atomic.Bool
)

func main() {
	// Wait for USB serial to start
	for !machine.USBCDC.DTR() {
		time.Sleep(time.Millisecond * 10)
	}

	// Init radio
	spi.Configure(machine.SPIConfig{})

	radio = sx127x.New(spi, rst, cs, d0) //machine.NoPin)

	for {
		cmdData := readCommand()
		cmd := strings.Split(string(cmdData), " ")

		switch string(cmd[0]) {
		case "sniff":
			// wait for any outstanding radio activity to cease
			stopRadio.Store(true)
			radioWaitGroup.Wait()

			ch := 15
			if len(cmd) > 1 {
				chParam, err := strconv.Atoi(cmd[1])
				if err == nil {
					ch = chParam
				}
			}

			stopRadio.Store(false)
			radioWaitGroup.Add(1)
			go sniff(ch)
		case "reset":
			stopRadio.Store(true)
			radioWaitGroup.Wait()
		default:
			println("unknown command")
		}
	}
}

func sniff(ch int) {
	defer radioWaitGroup.Done()

	//430.1 MHz + ch * 200 kHz
	freq := uint32(430_100_000) + uint32(ch)*200_000

	err := radio.Configure(sx127x.Config{
		Frequency:       freq,
		SpreadingFactor: 7,
		Bandwidth:       125000,
		CRC:             sx127x.CrcModeOn,
	})
	if err != nil {
		panic(err)
	}

	if !radio.Detect() {
		panic("radio not detected")
	}

	for {
		buf, err := radio.Rx(1000)
		if err != nil {
			panic(err)
		}

		if len(buf) > 2 && buf[0] == 0x88 && buf[1] == 0x41 {
			machine.Serial.Write([]byte{byte(len(buf) >> 8), byte(len(buf))})
			machine.Serial.Write(buf)
		}

		if stopRadio.Load() {
			return
		}
	}

}

func readCommand() []byte {
	buf := make([]byte, 0, 20)

	for {
		b, err := machine.Serial.ReadByte()
		if err != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		if b == '\r' || b == '\n' {
			return buf
		}

		buf = append(buf, b)
	}
}
