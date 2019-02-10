package znp

import (
	"errors"
	"fmt"
	"time"

	unp "github.com/dyrkin/unp-go"

	"github.com/dyrkin/bin"
	"github.com/dyrkin/znp-go/reflection"
)

type Sync struct {
	frame   *unp.Frame
	syncRsp chan *unp.Frame
	syncErr chan error
}

type Async struct {
	frame *unp.Frame
}

type Outgoing interface {
	Frame() *unp.Frame
}

func (s *Sync) Frame() *unp.Frame {
	return s.frame
}

func (a *Async) Frame() *unp.Frame {
	return a.frame
}

type Znp struct {
	u            *unp.Unp
	outbound     chan Outgoing
	inbound      chan *unp.Frame
	AsyncInbound chan interface{}
	Errors       chan error
	InFramesLog  chan *unp.Frame
	OutFramesLog chan *unp.Frame
	logInFrames  bool
	logOutFrames bool
}

func New(u *unp.Unp) *Znp {
	znp := &Znp{
		u:            u,
		outbound:     make(chan Outgoing),
		inbound:      make(chan *unp.Frame),
		AsyncInbound: make(chan interface{}),
		Errors:       make(chan error),
		InFramesLog:  make(chan *unp.Frame),
		OutFramesLog: make(chan *unp.Frame),
	}
	znp.startProcessor()
	go znp.incomingLoop()
	return znp
}

func (znp *Znp) LogInFrames(enabled bool) {
	znp.logInFrames = enabled
}

func (znp *Znp) LogOutFrames(enabled bool) {
	znp.logOutFrames = enabled
}

func (znp *Znp) ProcessRequest(commandType unp.CommandType, subsystem unp.Subsystem, command byte, request interface{}, response interface{}) (err error) {
	frame := &unp.Frame{
		CommandType: commandType,
		Subsystem:   subsystem,
		Command:     command,
		Payload:     bin.Encode(request),
	}
	done := make(chan bool, 1)
	go func() {
		if commandType == unp.C_SREQ {
			outgoing := &Sync{frame: frame,
				syncRsp: make(chan *unp.Frame, 1),
				syncErr: make(chan error, 1),
			}
			znp.outbound <- outgoing
			select {
			case frame := <-outgoing.syncRsp:
				bin.Decode(frame.Payload, response)
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
	syncRsp := make(chan *unp.Frame)
	syncErr := make(chan error)
	outgoingProcessor := func() {
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
					znp.u.WriteFrame(req.frame)
					responseProcessor := func() {
						select {
						case _ = <-deadline.timer.C:
							if !deadline.timer.Stop() {
								req.syncErr <- fmt.Errorf("timed out while waiting response for command: 0x%x sent to subsystem: %s ", frame.Command, frame.Subsystem)
							}
						case response := <-syncRsp:
							deadline.Cancel()
							req.syncRsp <- response
						case err := <-syncErr:
							deadline.Cancel()
							req.syncErr <- err
						}
					}
					go responseProcessor()
					logFrame(frame, znp.logOutFrames, znp.OutFramesLog)
				case *Async:
					znp.u.WriteFrame(req.frame)
					logFrame(req.frame, znp.logOutFrames, znp.OutFramesLog)
				}
			}
		}
	}
	incomingProcessor := func() {
		for {
			select {
			case frame := <-znp.inbound:
				switch frame.CommandType {
				case unp.C_SRSP:
					if frame.Subsystem == unp.S_RES0 && frame.Command == 0 {
						errorCode := frame.Payload[0]
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
						syncErr <- errors.New(errorMessage)
					} else {
						syncRsp <- frame
					}
				case unp.C_AREQ:
					key := registryKey{frame.Subsystem, frame.Command}
					if value, ok := asyncCommandRegistry[key]; ok {
						copy := reflection.Copy(value)
						bin.Decode(frame.Payload, copy)
						znp.AsyncInbound <- copy
					} else {
						znp.Errors <- fmt.Errorf("Unknown async command received: %v", frame)
					}
				}
			}
		}
	}
	go incomingProcessor()
	go outgoingProcessor()
}

func (znp *Znp) incomingLoop() {
	for {
		frame, err := znp.u.ReadFrame()
		if err != nil {
			znp.Errors <- err
		} else {
			logFrame(frame, znp.logInFrames, znp.InFramesLog)
			znp.inbound <- frame
		}
	}
}

func logFrame(frame *unp.Frame, log bool, logger chan *unp.Frame) {
	if log {
		logger <- frame
	}
}
