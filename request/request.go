package request

import unp "github.com/dyrkin/unp-go"

type Outgoing interface {
	Frame() *unp.Frame
}
