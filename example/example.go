package main

import (
	"fmt"
	"log"
	"time"

	unpi "github.com/dyrkin/unpi-go"
	znp "github.com/dyrkin/znp-go"
	serial "go.bug.st/serial.v1"
)

func main() {
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
				fmt.Printf("Async: %v\n", async)
			case frame := <-z.FramesLog:
				fmt.Printf("Frame received: %v\n", frame)
			}
		}
	}()

	// z.SysResetReq(1)

	var res interface{}

	res, err = z.SysPing()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", res)

	res, err = z.SysVersion()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", res)

	res, err = z.SysSetExtAddr("0x00124b00019c2ee9")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", res)

	res, err = z.SysGetExtAddr()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", res)

	res, err = z.SapiZbStartRequest()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", res)

	res, err = z.SapiZbPermitJoiningRequest("0xFF00", 200)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", res)

	res, err = z.LedControl(1, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", res)

	res, err = z.SapiZbReadConfiguration(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", res)

	res, err = z.SapiZbFindDeviceRequest("0x00124b00019c2ee9")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", res)

	res, err = z.SysOsalStartTimer(1, 3000)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", res)

	t := time.Now()

	res, err = z.SysSetTime(0, uint8(t.Hour()), uint8(t.Minute()), uint8(t.Second()), uint8(t.Month()), uint8(t.Day()), uint16(t.Year()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", res)

	res, err = z.SysGetTime()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", res)

	time.Sleep(200 * time.Second)
}
