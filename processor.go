package znp

import (
	"errors"
	"fmt"
	"time"

	"github.com/dyrkin/unp-go"

	"github.com/dyrkin/bin"
	"github.com/dyrkin/znp-go/reflection"
	"github.com/dyrkin/znp-go/request"
)

func (znp *Znp) Start() {
	startProcessors(znp)
	startIncomingFrameLoop(znp)
	znp.started = true
}

func (znp *Znp) Stop() {
	znp.started = false
}

func (znp *Znp) ProcessRequest(commandType unp.CommandType, subsystem unp.Subsystem, command byte, req interface{}, resp interface{}) error {
	frame := &unp.Frame{
		CommandType: commandType,
		Subsystem:   subsystem,
		Command:     command,
		Payload:     bin.Encode(req),
	}
	return processFrame(znp, frame, resp)
}

func processFrame(znp *Znp, frame *unp.Frame, resp interface{}) (err error) {
	if !znp.started {
		panic("Znp is not started. Call znp.Start() before")
	}
	completed := make(chan bool, 1)
	go func() {
		switch frame.CommandType {
		case unp.C_SREQ:
			outgoing := request.NewSync(frame)
			znp.outbound <- outgoing
			select {
			case frame := <-outgoing.SyncRsp():
				bin.Decode(frame.Payload, resp)
			case err = <-outgoing.SyncErr():
			}
		case unp.C_AREQ:
			outgoing := request.NewAsync(frame)
			znp.outbound <- outgoing
		default:
			err = fmt.Errorf("unsupported command type: %s ", frame.CommandType)
		}
		completed <- true
	}()
	<-completed
	return
}

func startProcessors(znp *Znp) {
	syncRsp := make(chan *unp.Frame)
	syncErr := make(chan error)
	syncRequestProcessor := makeSyncRequestProcessor(znp, syncRsp, syncErr)
	asyncRequestProcessor := makeAsyncRequestProcessor(znp)
	syncResponseProcessor := makeSyncResponseProcessor(syncRsp, syncErr)
	asyncResponseProcessor := makeAsyncResponseProcessor(znp)
	outgoingProcessor := func() {
		for znp.started {
			select {
			case outgoing := <-znp.outbound:
				switch req := outgoing.(type) {
				case *request.Sync:
					syncRequestProcessor(req)
				case *request.Async:
					asyncRequestProcessor(req)
				}
			}
		}
	}
	incomingProcessor := func() {
		for znp.started {
			select {
			case frame := <-znp.inbound:
				switch frame.CommandType {
				case unp.C_SRSP:
					syncResponseProcessor(frame)
				case unp.C_AREQ:
					asyncResponseProcessor(frame)
				default:
					select {
					case znp.errors <- fmt.Errorf("unsupported frame received type: %v ", frame):
					default:
					}
				}
			}
		}
	}
	go incomingProcessor()
	go outgoingProcessor()
}

func startIncomingFrameLoop(znp *Znp) {
	incomingLoop := func() {
		for znp.started {
			frame, err := znp.u.ReadFrame()
			if err != nil {
				select {
				case znp.errors <- err:
				default:
				}
			} else {
				logFrame(frame, znp.inFramesLog)
				znp.inbound <- frame
			}
		}
	}
	go incomingLoop()
}

func makeSyncRequestProcessor(znp *Znp, syncRsp chan *unp.Frame, syncErr chan error) func(req *request.Sync) {
	return func(req *request.Sync) {
		frame := req.Frame()
		deadline := time.NewTimer(5 * time.Second)
		logFrame(frame, znp.outFramesLog)
		err := znp.u.WriteFrame(frame)
		if err != nil {
			req.SyncErr() <- err
			return
		}
		select {
		case _ = <-deadline.C:
			if !deadline.Stop() {
				req.SyncErr() <- fmt.Errorf("timed out while waiting response for command: 0x%x sent to subsystem: %s ", frame.Command, frame.Subsystem)
			}
		case response := <-syncRsp:
			deadline.Stop()
			req.SyncRsp() <- response
		case err := <-syncErr:
			deadline.Stop()
			req.SyncErr() <- err
		}
	}
}

func makeAsyncRequestProcessor(znp *Znp) func(req *request.Async) {
	return func(req *request.Async) {
		logFrame(req.Frame(), znp.outFramesLog)
		znp.u.WriteFrame(req.Frame())
	}
}

var errorMessages = map[uint8]string{
	1: "Invalid subsystem",
	2: "Invalid command ID",
	3: "Invalid parameter",
	4: "Invalid length",
}

func makeSyncResponseProcessor(syncRsp chan *unp.Frame, syncErr chan error) func(frame *unp.Frame) {
	return func(frame *unp.Frame) {
		if frame.Subsystem == unp.S_RES0 && frame.Command == 0 {
			errorCode := frame.Payload[0]
			syncErr <- errors.New(errorMessages[errorCode])
		} else {
			syncRsp <- frame
		}
	}
}

func makeAsyncResponseProcessor(znp *Znp) func(frame *unp.Frame) {
	return func(frame *unp.Frame) {
		key := key{frame.Subsystem, frame.Command}
		if value, ok := asyncCommandRegistry[key]; ok {
			cp := reflection.Copy(value)
			bin.Decode(frame.Payload, cp)
			select {
			case znp.asyncInbound <- cp:
			default:
			}
		} else {
			select {
			case znp.errors <- fmt.Errorf("unknown async command received: %v", frame):
			default:
			}
		}
	}
}

func logFrame(frame *unp.Frame, logger chan *unp.Frame) {
	go func() {
		select {
		case logger <- frame:
		default:
		}
	}()

}
