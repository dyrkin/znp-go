package util

import (
	"errors"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestUintToHexString(c *C) {
	var str interface{}
	str, _ = UintToHexString(uint8(1))
	c.Assert(str, Equals, "0x01")
	str, _ = UintToHexString(uint16(1))
	c.Assert(str, Equals, "0x0001")
	str, _ = UintToHexString(uint32(1))
	c.Assert(str, Equals, "0x00000001")
	str, _ = UintToHexString(uint64(1))
	c.Assert(str, Equals, "0x0000000000000001")
	_, err := UintToHexString(int(1))
	c.Assert(err, DeepEquals, errors.New("Unsupported value: %!s(int=1)"))
}
