package main

import (
	"fmt"
	"log"

	cc "github.com/dyrkin/cc-go"
	unpi "github.com/dyrkin/unpi-go"
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
	c := cc.New(u)
	go func() {
		for {
			select {
			case err := <-c.Errors:
				fmt.Printf("Error: %s", err)
			case async := <-c.AsyncInbound:
				fmt.Printf("Async: %v", async)
			}
		}
	}()
	ping, err := c.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", ping)

	version, err := c.Version()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", version)

	enabledLed, err := c.LedControl(1, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", enabledLed)
}
