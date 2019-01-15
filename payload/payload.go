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

type tags struct {
	hex        tag
	endianness tag
	size       tag
	bitmask    tag
	bits       tag
	bound      tag
}

var emptyTag = &tags{}

func newTags(field reflect.StructField) *tags {
	t := field.Tag
	if t != "" {
		hex := tag(t.Get("hex"))
		endianness := tag(t.Get("endianness"))
		size := tag(t.Get("size"))
		bitmask := tag(t.Get("bitmask"))
		bits := tag(t.Get("bits"))
		bound := tag(t.Get("bound"))
		return &tags{hex: hex,
			endianness: endianness,
			size:       size,
			bitmask:    bitmask,
			bits:       bits,
			bound:      bound,
		}
	} else {
		return emptyTag
	}
}

func convertTo(v interface{}, typ reflect.Type) interface{} {
	value := reflect.ValueOf(v)
	return valueConvertTo(value, typ).Interface()
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
