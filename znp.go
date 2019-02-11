package znp

import (
	unp "github.com/dyrkin/unp-go"

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
	logInFrames  bool
	logOutFrames bool
	started      bool
}

func New(u *unp.Unp) *Znp {
	znp := &Znp{
		u:            u,
		outbound:     make(chan request.Outgoing),
		inbound:      make(chan *unp.Frame),
		asyncInbound: make(chan interface{}),
		errors:       make(chan error),
		inFramesLog:  make(chan *unp.Frame, 100),
		outFramesLog: make(chan *unp.Frame, 100),
	}
	return znp
}

func (znp *Znp) LogInFrames(enabled bool) {
	znp.logInFrames = enabled
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

func (znp *Znp) LogOutFrames(enabled bool) {
	znp.logOutFrames = enabled
}

func (znp *Znp) IsStarted() bool {
	return znp.started
}
