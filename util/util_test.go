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
	str, _ = UintToHexString(uint64(1), 1)
	c.Assert(str, Equals, "0x01")
	str, _ = UintToHexString(uint64(1), 2)
	c.Assert(str, Equals, "0x0001")
	str, _ = UintToHexString(uint64(1), 4)
	c.Assert(str, Equals, "0x00000001")
	str, _ = UintToHexString(uint64(1), 8)
	c.Assert(str, Equals, "0x0000000000000001")
	_, err := UintToHexString(uint64(1), 7)
	c.Assert(err, DeepEquals, errors.New("Unsupported size: 7"))
}
