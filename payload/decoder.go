package payload

import (
	"bytes"
	"encoding/binary"

	"reflect"
	"strconv"

	"github.com/dyrkin/znp-go/util"
)

type decoder struct {
	buf *bytes.Buffer
}

//Decode struct from byte array
func Decode(payload []byte, response interface{}) {
	value := reflect.ValueOf(response)
	decoder := &decoder{bytes.NewBuffer(payload)}
	decoder.decode(value)
}

func (d *decoder) decode(value reflect.Value) {
	switch value.Kind() {
	case reflect.Ptr:
		d.pointer(value)
	case reflect.Struct:
		d.strukt(value)
	}
}

func (d *decoder) pointer(value reflect.Value) {
	if value.IsNil() {
		element := reflect.New(value.Type().Elem())
		if value.CanSet() {
			value.Set(element)
		}
	}
	d.decode(value.Elem())
}

func (d *decoder) strukt(value reflect.Value) {
	var bitmaskBytes uint64
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := value.Type().Field(i)
		tags := newTags(fieldType)
		switch field.Kind() {
		case reflect.Ptr:
			d.pointer(field)
		case reflect.String:
			d.string(field, tags)
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			d.uint(field, tags, &bitmaskBytes)
		case reflect.Array:
			d.array(field, tags)
		case reflect.Slice:
			d.slice(field, tags)
		}
	}
}

func (d *decoder) slice(value reflect.Value, tags *tags) {
	length := d.dynamicLength(tags)
	value.Set(reflect.MakeSlice(value.Type(), length, length))
	for i := 0; i < length; i++ {
		sliceElement := value.Index(i)
		switch sliceElement.Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			d.uint(sliceElement, tags, nil)
		case reflect.String:
			d.string(sliceElement, tags)
		case reflect.Ptr:
			d.pointer(sliceElement)
		case reflect.Struct:
			d.strukt(sliceElement)
		}
	}
}

func (d *decoder) array(value reflect.Value, tags *tags) {
	if value.Len() > 0 {
		size := int(value.Index(0).Type().Size())
		for i := 0; i < value.Len(); i++ {
			arrayElem := value.Index(i)
			v := d.readUint(tags.endianness, size)
			arrayElem.SetUint(v)
		}
	}
}

func (d *decoder) uint(value reflect.Value, tags *tags, bitmaskBytes *uint64) {
	if value.CanAddr() {
		if tags.bits.nonEmpty() {
			if tags.bitmask == "start" {
				*bitmaskBytes = d.readUint(tags.endianness, int(value.Type().Size()))
			}
			bitmaskBits := bitmaskBits(tags.bits)
			pos := util.FirstBitPosition(bitmaskBits)
			v := (*bitmaskBytes & bitmaskBits) >> pos
			value.SetUint(v)
		} else if tags.bound.nonEmpty() {
			size, _ := strconv.Atoi(string(tags.bound))
			v := d.readUint(tags.endianness, size)
			value.SetUint(v)
		} else {
			v := d.readUint(tags.endianness, int(value.Type().Size()))
			value.SetUint(v)
		}
	} else {
		panic("Unaddressable uint value")
	}
}

func (d *decoder) string(value reflect.Value, tags *tags) {
	if tags.hex.nonEmpty() {
		size, _ := strconv.Atoi(string(tags.hex))
		v := d.readUint(tags.endianness, size)
		hexString, _ := util.UintToHexString(v, size)
		value.SetString(hexString)
	} else {
		length := d.dynamicLength(tags)
		b := make([]uint8, length, length)
		d.read(tags.endianness, b)
		value.SetString(string(b))
	}
}

func (d *decoder) dynamicLength(tags *tags) int {
	if tags.size.nonEmpty() {
		size, _ := strconv.Atoi(string(tags.size))
		return int(d.readUint(tags.endianness, size))
	} else {
		return len(d.buf.Bytes())
	}
}

func (d *decoder) read(endianness tag, v interface{}) {
	binary.Read(d.buf, order(endianness), v)
}

func (d *decoder) read2(endianness tag, v reflect.Value) {
	switch v.Kind() {
	case reflect.Uint8:
		u, _ := d.buf.ReadByte()
		v.SetUint(uint64(u))
	case reflect.Uint16:
		u := d.readUint(endianness, 2)
		v.SetUint(u)
	case reflect.Uint32:
		u := d.readUint(endianness, 4)
		v.SetUint(u)
	case reflect.Uint64:
		u := d.readUint(endianness, 8)
		v.SetUint(u)
	}
}

func (d *decoder) readUint(endianness tag, size int) uint64 {
	var v uint64
	if endianness == "be" {
		for i := 0; i < size; i++ {
			t, _ := d.buf.ReadByte()
			v = v | uint64(t)<<byte((size-i-1)*8)
		}
	} else {
		for i := 0; i < size; i++ {
			t, _ := d.buf.ReadByte()
			v = v | uint64(t)<<byte(i*8)
		}
	}
	return v
}
