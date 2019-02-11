package request

import unp "github.com/dyrkin/unp-go"

type Sync struct {
	frame   *unp.Frame
	syncRsp chan *unp.Frame
	syncErr chan error
}

func NewSync(frame *unp.Frame) *Sync {
	return &Sync{frame, make(chan *unp.Frame, 1), make(chan error, 1)}
}

func (s *Sync) Frame() *unp.Frame {
	return s.frame
}

func (s *Sync) SyncRsp() chan *unp.Frame {
	return s.syncRsp
}

func (s *Sync) SyncErr() chan error {
	return s.syncErr
}
