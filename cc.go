package cc

import (
	"fmt"
	"time"

	unpi "github.com/dyrkin/unpi-go"
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

type Cc struct {
	u            *unpi.Unpi
	outbound     chan Outgoing
	inbound      chan *unpi.Frame
	AsyncInbound chan *unpi.Frame
	Errors       chan error
}

type ResetRequest struct {
	ResetType byte
}

type PingResponse struct {
	Capabilities uint16
}

type VersionResponse struct {
	Transportrev uint8
	Product      uint8
	Majorrel     uint8
	Minorrel     uint8
	Maintrel     uint8
}

type LedControlRequest struct {
	LedID uint8
	Mode  uint8
}

type LedControlResponse struct {
	Status uint8
}

func New(u *unpi.Unpi) *Cc {
	cc := &Cc{
		u:            u,
		outbound:     make(chan Outgoing),
		inbound:      make(chan *unpi.Frame),
		AsyncInbound: make(chan *unpi.Frame),
		Errors:       make(chan error),
	}
	go cc.startProcessor()
	go cc.incomingLoop()
	return cc
}

func (cc *Cc) Reset(resetType byte) {
	cc.ProcessRequest(unpi.C_AREQ, unpi.S_SYS, 0, &ResetRequest{resetType}, nil)
}

func (cc *Cc) Ping() (*PingResponse, error) {
	rsp := &PingResponse{}
	err := cc.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 1, nil, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (cc *Cc) Version() (*VersionResponse, error) {
	rsp := &VersionResponse{}
	err := cc.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 2, nil, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (cc *Cc) LedControl(ledid uint8, mode uint8) (*LedControlResponse, error) {
	req := &LedControlRequest{LedID: ledid, Mode: mode}
	rsp := &LedControlResponse{}
	err := cc.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 10, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (cc *Cc) ProcessRequest(commandType unpi.CommandType, subsystem unpi.Subsystem, command byte, request interface{}, response interface{}) (err error) {
	frame := &unpi.Frame{
		CommandType: commandType,
		Subsystem:   subsystem,
		Command:     command,
		Payload:     serialize(request),
	}
	done := make(chan bool, 1)
	go func() {
		if commandType == unpi.C_SREQ {
			outgoing := &Sync{frame: frame,
				syncRsp: make(chan *unpi.Frame, 1),
				syncErr: make(chan error, 1),
			}
			cc.outbound <- outgoing
			select {
			case frame := <-outgoing.syncRsp:
				deserialize(frame.Payload, response)
			case err = <-outgoing.syncErr:
			}
		} else {
			outgoing := &Async{frame: frame}
			cc.outbound <- outgoing
		}
		done <- true
	}()
	<-done
	return
}

func (cc *Cc) startProcessor() {
	registry := NewRequestRegistry()
	for {
		select {
		case outgoing := <-cc.outbound:
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
				cc.u.Write(req.frame)
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
				cc.u.Write(req.frame)
			}
		case frame := <-cc.inbound:
			if frame.CommandType == unpi.C_SRSP {
				key := &registryKey{frame.Subsystem, frame.Command}
				value, ok := registry.Get(key)
				if !ok {
					cc.Errors <- fmt.Errorf("Unknown response received: %v", frame)
					continue
				}
				value.deadline.Cancel()
				value.syncRsp <- frame
			} else {
				cc.AsyncInbound <- frame
			}
		}
	}
}

func (cc *Cc) incomingLoop() {
	for {
		frame, err := cc.u.Read()
		if err != nil {
			cc.Errors <- err
		} else {
			cc.inbound <- frame
		}
	}
}
