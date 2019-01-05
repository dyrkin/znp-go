package main

import (
	"fmt"
	"log"

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
	go func() {
		for {
			select {
			case err := <-z.Errors:
				fmt.Printf("Error: %s", err)
			case async := <-z.AsyncInbound:
				fmt.Printf("Async: %v", async)
			}
		}
	}()

	// c.Reset(1)

	ping, err := z.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", ping)

	version, err := z.Version()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", version)

	enabledLed, err := z.LedControl(1, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", enabledLed)
}
