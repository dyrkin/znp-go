# ZigBee Network Processor (ZNP) Interface

[![Build Status](https://cloud.drone.io/api/badges/dyrkin/znp-go/status.svg??branch=master)](https://cloud.drone.io/dyrkin/znp-go)

## Overview

ZNP is used for communication between the host and a ZigBee device through a serial port. You can issue Monitor and Test (MT) commands to the
ZigBee target from your application.

I tested it with cc253**1**, but it might work with cc253**X**  

## Example

To use it you need to provide a reference to an unp instance:

```go
import (
	"go.bug.st/serial.v1"
	"github.com/davecgh/go-spew/spew"
	"github.com/dyrkin/unp-go"
	"github.com/dyrkin/znp-go"
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

	u := unp.New(1, port)
	z := znp.New(u)
	z.Start()
}
```

Then you be able to run commands:

```go
res, err := z.SysSetExtAddr("0x00124b00019c2ee9")
if err != nil {
	log.Fatal(err)
}
	
res, err = z.SapiZbPermitJoiningRequest("0xFF00", 200)
if err != nil {
	log.Fatal(err)
}
```

To receive async commands and errors, use `AsyncInbound()` and `Errors()` channels:

```go
go func() {
    for {
        select {
        case err := <-z.Errors():
            fmt.Printf("Error received: %s\n", err)
        case async := <-z.AsyncInbound():
            fmt.Printf("Async received: %s\n", spew.Sdump(async))
        }
    }
}()
```

To log all ingoing and outgoing unp frames, use `InFramesLog()` and `OutFramesLog()` channels:

```go
go func() {
    for {
        select {
        case frame := <-z.OutFramesLog():
            fmt.Printf("Frame sent: %s\n", spew.Sdump(frame))
        case frame := <-z.InFramesLog():
            fmt.Printf("Frame received: %s\n", spew.Sdump(frame))
        }
    }
}()
```

See more [examples](example/example.go)

