package payload

import (
	"encoding/binary"
	"reflect"
	"strconv"
)

var types = map[int]reflect.Type{
	1: reflect.TypeOf(uint8(0)),
	2: reflect.TypeOf(uint16(0)),
	4: reflect.TypeOf(uint32(0)),
	8: reflect.TypeOf(uint64(0)),
}

var intType = reflect.TypeOf(int(0))
var uint64Type = reflect.TypeOf(uint64(0))

type tag string

func (t tag) nonEmpty() bool {
	return len(t) > 0
}

type tags reflect.StructTag

func (t tags) hex() tag {
	return tag(reflect.StructTag(t).Get("hex"))
}

func (t tags) endianness() tag {
	return tag(reflect.StructTag(t).Get("endianness"))
}

func (t tags) size() tag {
	return tag(reflect.StructTag(t).Get("size"))
}

func (t tags) bitmask() tag {
	return tag(reflect.StructTag(t).Get("bitmask"))
}

func (t tags) bits() tag {
	return tag(reflect.StructTag(t).Get("bits"))
}

func (t tags) bound() tag {
	return tag(reflect.StructTag(t).Get("bound"))
}

func valueConvertTo(value reflect.Value, typ reflect.Type) reflect.Value {
	return value.Convert(typ)
}

func bitmaskBits(value tag) (bitmaskBits uint64) {
	prefix := string(value[:2])
	bitmask := string(value[2:])
	if prefix == "0x" {
		bitmaskBits, _ = strconv.ParseUint(bitmask, 16, len(bitmask)*4)
		return
	} else if prefix == "0b" {
		bitmaskBits, _ = strconv.ParseUint(bitmask, 2, len(bitmask))
		return
	}
	panic("Unsupported prefix: " + prefix)
}

func order(endianness tag) binary.ByteOrder {
	if endianness == "be" {
		return binary.BigEndian
	}
	return binary.LittleEndian
}
