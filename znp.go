package znp

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/sys/unix"

	unpi "github.com/dyrkin/unpi-go"

	"github.com/dyrkin/znp-go/payload"
	"github.com/dyrkin/znp-go/reflection"
)

type Sync struct {
	frame   *unpi.Frame
	syncRsp chan *unpi.Frame
	syncErr chan error
}

type Async struct {
	frame *unpi.Frame
}

type Outgoing interface {
	Frame() *unpi.Frame
}

func (s *Sync) Frame() *unpi.Frame {
	return s.frame
}

func (a *Async) Frame() *unpi.Frame {
	return a.frame
}

type Znp struct {
	u            *unpi.Unpi
	outbound     chan Outgoing
	inbound      chan *unpi.Frame
	AsyncInbound chan interface{}
	Errors       chan error
	FramesLog    chan *unpi.Frame
	logFrames    bool
}

func New(u *unpi.Unpi) *Znp {
	znp := &Znp{
		u:            u,
		outbound:     make(chan Outgoing),
		inbound:      make(chan *unpi.Frame),
		AsyncInbound: make(chan interface{}),
		Errors:       make(chan error),
		FramesLog:    make(chan *unpi.Frame),
	}
	go znp.startProcessor()
	go znp.incomingLoop()
	return znp
}

func (znp *Znp) LogFrames(enabled bool) {
	znp.logFrames = enabled
}

func (znp *Znp) ProcessRequest(commandType unpi.CommandType, subsystem unpi.Subsystem, command byte, request interface{}, response interface{}) (err error) {
	frame := &unpi.Frame{
		CommandType: commandType,
		Subsystem:   subsystem,
		Command:     command,
		Payload:     payload.Encode(request),
	}
	done := make(chan bool, 1)
	go func() {
		if commandType == unpi.C_SREQ {
			outgoing := &Sync{frame: frame,
				syncRsp: make(chan *unpi.Frame, 1),
				syncErr: make(chan error, 1),
			}
			znp.outbound <- outgoing
			select {
			case frame := <-outgoing.syncRsp:
				payload.Decode(frame.Payload, response)
			case err = <-outgoing.syncErr:
			}
		} else {
			outgoing := &Async{frame: frame}
			znp.outbound <- outgoing
		}
		done <- true
	}()
	<-done
	return
}

func (znp *Znp) startProcessor() {
	registry := NewRequestRegistry()
	for {
		select {
		case outgoing := <-znp.outbound:
			switch req := outgoing.(type) {
			case *Sync:
				frame := req.Frame()
				deadline := &deadline{
					time.NewTimer(5 * time.Second),
					make(chan bool, 1),
				}
				key := &registryKey{frame.Subsystem, frame.Command}
				value := &registryValue{req.syncRsp, req.syncErr, deadline}
				registry.Register(key, value)
				znp.u.WriteFrame(req.frame)
				go func() {
					select {
					case _ = <-deadline.timer.C:
						if !deadline.timer.Stop() {
							req.syncErr <- fmt.Errorf("timed out while waiting response for command: %b sent to subsystem: %s ", frame.Command, frame.Subsystem)
						}
					case _ = <-deadline.cancelled:
					}
					registry.Unregister(key)
				}()
			case *Async:
				znp.u.WriteFrame(req.frame)
			}
		case frame := <-znp.inbound:
			if frame.CommandType == unpi.C_SRSP {
				//process error response
				if frame.Subsystem == unpi.S_RES0 && frame.Command == 0 {
					errorCode := frame.Payload[0]
					subsystem := unpi.Subsystem(frame.Payload[1] & 0x1F)
					command := frame.Payload[2]
					key := &registryKey{subsystem, command}
					value, ok := registry.Get(key)
					if !ok {
						znp.Errors <- fmt.Errorf("Unknown response received: %v", frame)
						continue
					}
					value.deadline.Cancel()
					var errorMessage string
					switch errorCode {
					case 1:
						errorMessage = "Invalid subsystem"
					case 2:
						errorMessage = "Invalid command ID"
					case 3:
						errorMessage = "Invalid parameter"
					case 4:
						errorMessage = "Invalid length"
					}
					value.syncErr <- errors.New(errorMessage)
				} else {
					key := &registryKey{frame.Subsystem, frame.Command}
					value, ok := registry.Get(key)
					if !ok {
						znp.Errors <- fmt.Errorf("Unknown response received: %v", frame)
						continue
					}
					value.deadline.Cancel()
					value.syncRsp <- frame
				}
			} else {
				key := registryKey{frame.Subsystem, frame.Command}
				if value, ok := AsyncCommandRegistry[key]; ok {
					copy := reflection.Copy(value)
					payload.Decode(frame.Payload, copy)
					znp.AsyncInbound <- copy
				} else {
					znp.Errors <- fmt.Errorf("Unknown async command received: %v", frame)
				}
			}
		}
	}
}

func (znp *Znp) incomingLoop() {
	for {
		frame, err := znp.u.ReadFrame()
		if err != nil {
			znp.Errors <- err
			if err == unix.ENXIO {
				time.Sleep(5 * time.Second)
			}
		} else {
			znp.inbound <- frame
			if znp.logFrames {
				znp.FramesLog <- frame
			}
		}
	}
}
