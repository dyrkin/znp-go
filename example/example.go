package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	unpi "github.com/dyrkin/unpi-go"
	znp "github.com/dyrkin/znp-go"
	serial "go.bug.st/serial.v1"
)

func main() {
	spew.Config.DisableCapacities = true
	spew.Config.DisablePointerAddresses = true
	mode := &serial.Mode{
		BaudRate: 115200,
	}

	port, err := serial.Open("/dev/tty.usbmodem14101", mode)
	if err != nil {
		log.Fatal(err)
	}
	port.SetRTS(true)

	u := unpi.New(1, port)
	z := znp.New(u)
	z.LogFrames(false)
	go func() {
		for {
			select {
			case err := <-z.Errors:
				fmt.Printf("Error: %s\n", err)
			case async := <-z.AsyncInbound:
				fmt.Printf("Async: %s\n", spew.Sdump(async))
			case frame := <-z.FramesLog:
				fmt.Printf("Frame received: %s\n", spew.Sdump(frame))
			}
		}
	}()

	// z.SysResetReq(1)

	var res interface{}

	res, err = z.SysPing()
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SysVersion()
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SysSetExtAddr("0x00124b00019c2ee9")
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SysGetExtAddr()
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SapiZbStartRequest()
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SapiZbPermitJoiningRequest("0xFF00", 200)
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.UtilLedControl(1, znp.OFF)
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SapiZbReadConfiguration(1)
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SapiZbFindDeviceRequest("0x00124b00019c2ee9")
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SysOsalStartTimer(1, 3000)
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	t := time.Now()

	res, err = z.SysSetTime(0, uint8(t.Hour()), uint8(t.Minute()), uint8(t.Second()), uint8(t.Month()), uint8(t.Day()), uint16(t.Year()))
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SysGetTime()
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.UtilGetDeviceInfo()
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.UtilGetNvInfo()
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.UtilLoopback([]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9})
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.UtilAssocFindDevice(1)
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.UtilAssocGetWithAddr("0x0000000000000000", "0x25cc")
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.UtilZclKeyEstSign([]uint8{1, 2})
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	time.Sleep(200 * time.Second)
}

func PrintStruct(v interface{}) {
	spew.Dump(v)
}

func RenderStruct(v interface{}) string {
	jsonBytes, _ := json.Marshal(v)
	return string(jsonBytes)
}
