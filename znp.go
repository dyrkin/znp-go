package znp

import (
	"github.com/dyrkin/unp-go"

	"github.com/dyrkin/znp-go/request"
)

type Znp struct {
	u            *unp.Unp
	outbound     chan request.Outgoing
	inbound      chan *unp.Frame
	asyncInbound chan interface{}
	errors       chan error
	inFramesLog  chan *unp.Frame
	outFramesLog chan *unp.Frame
	started      bool
}

func New(u *unp.Unp) *Znp {
	znp := &Znp{
		u:            u,
		outbound:     make(chan request.Outgoing),
		inbound:      make(chan *unp.Frame),
		asyncInbound: make(chan interface{}),
		errors:       make(chan error, 100),
		inFramesLog:  make(chan *unp.Frame, 100),
		outFramesLog: make(chan *unp.Frame, 100),
	}
	return znp
}

func (znp *Znp) Errors() chan error {
	return znp.errors
}

func (znp *Znp) AsyncInbound() chan interface{} {
	return znp.asyncInbound
}

func (znp *Znp) InFramesLog() chan *unp.Frame {
	return znp.inFramesLog
}

func (znp *Znp) OutFramesLog() chan *unp.Frame {
	return znp.outFramesLog
}

func (znp *Znp) IsStarted() bool {
	return znp.started
}
