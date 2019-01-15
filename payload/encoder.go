package payload

import (
	"bytes"
	"reflect"
	"strconv"

	"github.com/dyrkin/znp-go/util"
)

type encoder struct {
	buf *bytes.Buffer
}

//Encode struct to byte array
func Encode(request interface{}) []byte {
	value := reflect.ValueOf(request)
	buf := bytes.NewBuffer(make([]byte, 0, 200))
	encoder := &encoder{buf}
	encoder.encode(value)
	return buf.Bytes()
}

func (e *encoder) encode(value reflect.Value) {
	switch value.Kind() {
	case reflect.Ptr:
		e.pointer(value)
	case reflect.Struct:
		e.strukt(value)
	}
}

func (e *encoder) strukt(value reflect.Value) {
	var bitmaskBytes uint64
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := value.Type().Field(i)
		tags := newTags(fieldType)
		switch field.Kind() {
		case reflect.Ptr:
			e.pointer(field)
		case reflect.String:
			e.string(field, tags)
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			e.uint(field, tags, &bitmaskBytes)
		case reflect.Array:
			e.array(field, tags)
		case reflect.Slice:
			e.slice(field, tags)
		case reflect.Interface:
			e.pointer(field)
		}
	}
}

func (e *encoder) slice(value reflect.Value, tags *tags) {
	length := value.Len()
	e.dynamicLength(length, tags)
	for i := 0; i < length; i++ {
		sliceElement := value.Index(i)
		switch sliceElement.Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			e.uint(sliceElement, tags, nil)
		case reflect.String:
			e.string(sliceElement, tags)
		case reflect.Ptr:
			e.pointer(sliceElement)
		case reflect.Struct:
			e.strukt(sliceElement)
		}
	}
}

func (e *encoder) array(value reflect.Value, tags *tags) {
	for i := 0; i < value.Len(); i++ {
		e.write(tags.endianness, value.Index(i))
	}
}

func (e *encoder) string(value reflect.Value, tags *tags) {
	v := value.String()
	if tags.hex.nonEmpty() {
		size, _ := strconv.Atoi(string(tags.hex))
		typ := types[size]
		addr, _ := strconv.ParseUint(v[2:], 16, size*8)
		v := valueConvertTo(reflect.ValueOf(addr), typ)
		e.write(tags.endianness, v)
	} else {
		e.dynamicLength(len(v), tags)
		e.write(tags.endianness, reflect.ValueOf([]uint8(v)))
	}
}

func (e *encoder) uint(value reflect.Value, tags *tags, bitmaskBytes *uint64) {
	if tags.bits.nonEmpty() {
		if tags.bitmask == "start" {
			*bitmaskBytes = 0
		}
		bitmaskBits := bitmaskBits(tags.bits)
		pos := util.FirstBitPosition(bitmaskBits)
		v := valueConvertTo(value, uint64Type).Uint()
		v = ((v << pos) & bitmaskBits)
		*bitmaskBytes = (*bitmaskBytes) | v
		if tags.bitmask == "end" {
			v := valueConvertTo(reflect.ValueOf(*bitmaskBytes), value.Type())
			e.write(tags.endianness, v)
		}
	} else if tags.bound.nonEmpty() {
		size, _ := strconv.Atoi(string(tags.bound))
		e.writeUint(value.Uint(), size, tags.endianness)
	} else {
		e.write(tags.endianness, value)
	}
}

func (e *encoder) pointer(value reflect.Value) {
	e.encode(value.Elem())
}

func (e *encoder) dynamicLength(length int, tags *tags) {
	if tags.size.nonEmpty() {
		size, _ := strconv.Atoi(string(tags.size))
		if typ, ok := types[size]; ok {
			v := valueConvertTo(reflect.ValueOf(length), typ)
			e.write(tags.endianness, v)
		}
	}
}

func (e *encoder) write(endianness tag, v reflect.Value) {
	switch v.Kind() {
	case reflect.Uint8:
		e.buf.WriteByte(uint8(v.Uint()))
	case reflect.Uint16:
		e.writeUint(v.Uint(), 2, endianness)
	case reflect.Uint32:
		e.writeUint(v.Uint(), 4, endianness)
	case reflect.Uint64:
		e.writeUint(v.Uint(), 8, endianness)
	case reflect.Slice:
		l := v.Len()
		for i := 0; i < l; i++ {
			e.write(endianness, v.Index(i))
		}
	}
}

func (e *encoder) writeUint(t uint64, size int, endianness tag) {
	if endianness == "be" {
		for i := 0; i < size; i++ {
			e.buf.WriteByte(byte(t >> byte((size-i-1)*8)))
		}
	} else {
		for i := 0; i < size; i++ {
			e.buf.WriteByte(byte(t >> byte(i*8)))
		}
	}
}
