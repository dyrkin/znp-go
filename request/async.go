package request

import unp "github.com/dyrkin/unp-go"

type Async struct {
	frame *unp.Frame
}

func NewAsync(frame *unp.Frame) *Async {
	return &Async{frame}
}

func (a *Async) Frame() *unp.Frame {
	return a.frame
}
