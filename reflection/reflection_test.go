package reflection

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

type Struct struct {
	v1 uint8
	v2 string
}

func (s *MySuite) TestCopy(c *C) {
	copy1 := Copy(&Struct{1, "2"})

	c.Assert(copy1, DeepEquals, &Struct{})

	copy2 := Copy(Struct{1, "2"})

	c.Assert(copy2, DeepEquals, Struct{})
}
