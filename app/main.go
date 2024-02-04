package main

import (
	"encoding/hex"
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/kenbell/giv-lora/app/lora"
	"github.com/kenbell/giv-lora/app/modbus"
	"github.com/rivo/tview"
)

const (
	port    = "/dev/ttyACM0"
	channel = 26
)

var devices = map[byte]*modbus.Device{}

var (
	app        = tview.NewApplication()
	logView    = tview.NewTextView()
	errorView  = tview.NewTextView()
	meterViews = []*tview.TextView{
		tview.NewTextView(),
		tview.NewTextView(),
		tview.NewTextView()}
)

func main() {

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 {
			app.Stop()
		}

		return event
	})

	root := tview.NewFlex().SetDirection(tview.FlexRow)
	devices := tview.NewFlex()
	info := tview.NewFlex()

	for i, dv := range meterViews {
		dv.SetBorder(true)
		dv.SetTitle(fmt.Sprintf(" ID%d ", i+1))
		dv.SetWrap(false)

		devices.AddItem(dv, 0, 1, false)
	}

	logView.SetBorder(true)
	logView.SetBorderColor(tcell.ColorGreen)
	logView.SetTextColor(tcell.ColorGreen)
	logView.SetTitleColor(tcell.ColorGreen)
	logView.SetTitle(" log ")
	info.AddItem(logView, 0, 1, false)

	errorView.SetBorder(true)
	errorView.SetBorderColor(tcell.ColorRed)
	errorView.SetTextColor(tcell.ColorRed)
	errorView.SetTitleColor(tcell.ColorRed)
	errorView.SetTitle(" errors ")
	info.AddItem(errorView, 0, 1, false)

	root.AddItem(devices, 11, 1, false)
	root.AddItem(info, 0, 1, false)

	go handleStream(port)

	if err := app.SetRoot(root, true).Run(); err != nil {
		fmt.Printf("%v", err)
		time.Sleep(10 * time.Second)
		panic(err)
	}
}

func handleStream(port string) {
	defer func() {
		if r := recover(); r != nil {
			logError(fmt.Sprintf("%v", r))
			return
		}
	}()

	d, err := lora.NewDongle(port)
	if err != nil {
		logError(fmt.Sprintf("%v\n", err))
		return
	}

	err = d.StartSniff(channel)
	fmt.Fprintf(logView, "starting sniff 26\n")
	if err != nil {
		logError(fmt.Sprintf("%v\n", err))
		return
	}

	pktBuf := make([]byte, 256)

	for {
		n, err := d.Read(pktBuf)
		if err != nil {
			logError(fmt.Sprintf("%v\n", err))
			return
		}

		if n == 0 {
			continue
		}

		loraPkt := pktBuf[:n]

		seqNo, payload, err := lora.ParseLoRaPacket(loraPkt)
		if err != nil {
			fmt.Fprintln(errorView, "------------------")
			fmt.Fprintf(errorView, "LoRa Error: %v\n", err)
			fmt.Fprintln(errorView, hex.Dump(loraPkt))
			continue
		}

		pos := 0
		for pos < len(payload) {
			n, adu, err := modbus.ParseADU(payload[pos:])
			if err != nil {
				fmt.Fprintln(errorView, "------------------")
				fmt.Fprintf(errorView, "MODBUS Error: %v\n", err)
				fmt.Fprintf(errorView, "LoRa:\n")
				fmt.Fprintln(errorView, hex.Dump(loraPkt))
				fmt.Fprintf(errorView, "Modbus: %s\n", hex.EncodeToString(payload[pos:]))
				break
			}

			d, ok := devices[adu.Address]
			if !ok {
				d = modbus.NewDevice(adu.Address)
				devices[adu.Address] = d
			}

			err = d.HandleADU(adu)
			if err != nil {
				fmt.Fprintln(errorView, "------------------")
				fmt.Fprintf(errorView, "ADU Error: %v\n", err)
				fmt.Fprintf(errorView, "LoRa:\n")
				fmt.Fprintln(errorView, hex.Dump(loraPkt))
				fmt.Fprintf(errorView, "Modbus: %s\n", hex.EncodeToString(payload[pos:]))
				break
			}

			fmt.Fprintf(logView, "Seq: %d, %s\n", seqNo, adu.Dump())

			pos += n
		}

		for addr := byte(1); addr < 4; addr++ {
			updateDevice(addr)
		}

		logView.ScrollToEnd()
		errorView.ScrollToEnd()
		app.Draw()
	}

}

func updateDevice(addr byte) {
	tv := meterViews[addr-1]
	tv.Clear()

	d, ok := devices[addr]
	if !ok {
		return
	}

	keys := make([]uint16, 0, len(modbus.EnergyMeterRegisters))
	for k := range modbus.EnergyMeterRegisters {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for _, key := range keys {
		reg := key
		defn := modbus.EnergyMeterRegisters[key]

		dispValue := "--.--"

		v, ok := d.InputRegisters[reg]
		if ok {

			switch defn.DataType {
			case modbus.RegTypeFloat:
				v2, ok := d.InputRegisters[reg+1]
				if !ok {
					dispValue = fmt.Sprintf("err: 0x%04X", v)
					break
				}

				dispValue = fmt.Sprintf("%.2f", math.Float32frombits(uint32(v)<<16|uint32(v2)))

			default:
				dispValue = fmt.Sprintf("0x%04X", v)
			}
		}

		fmt.Fprintf(tv, "%20s: %s %s\n", defn.Name, dispValue, defn.Unit)
	}
}

func logError(msg string) {
	fmt.Fprintf(errorView, "%s: App Error: %s\n", time.Now().Format(time.ANSIC), msg)
}
