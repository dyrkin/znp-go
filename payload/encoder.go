package payload

import (
	"bytes"
	"encoding/binary"
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
	buf := &bytes.Buffer{}
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
	e.write(tags.endianness, value.Interface())
}

func (e *encoder) string(value reflect.Value, tags *tags) {
	v := value.Interface().(string)
	if tags.hex.nonEmpty() {
		size, _ := strconv.Atoi(string(tags.hex))
		typ := types[size]
		addr, _ := strconv.ParseUint(v[2:], 16, size*8)
		v := convertTo(addr, typ)
		e.write(tags.endianness, v)
	} else {
		e.dynamicLength(len(v), tags)
		e.write(tags.endianness, []uint8(v))
	}
}

func (e *encoder) uint(value reflect.Value, tags *tags, bitmaskBytes *uint64) {
	if tags.bits.nonEmpty() {
		if tags.bitmask == "start" {
			*bitmaskBytes = 0
		}
		bitmaskBits := bitmaskBits(tags.bits)
		pos := util.FirstBitPosition(bitmaskBits)
		v := valueConvertTo(value, uint64Type).Interface().(uint64)
		v = ((v << pos) & bitmaskBits)
		*bitmaskBytes = (*bitmaskBytes) | v
		if tags.bitmask == "end" {
			v := convertTo(*bitmaskBytes, value.Type())
			e.write(tags.endianness, v)
		}
	} else if tags.bound.nonEmpty() {
		size, _ := strconv.Atoi(string(tags.bound))
		v := make([]uint8, size, size)
		uintVal := value.Uint()
		endianness := tags.endianness
		for i := 0; i < size; i++ {
			if endianness == "be" {
				v[size-i-1] = byte(uintVal >> byte(i*8))
			} else {
				v[i] = byte(uintVal >> byte(i*8))
			}
		}
		e.write(tags.endianness, v)
	} else {
		v := value.Interface()
		e.write(tags.endianness, v)
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
			e.write(tags.endianness, v.Interface())
		}
	}
}

func (e *encoder) write(endianness tag, v interface{}) {
	binary.Write(e.buf, order(endianness), v)
}
