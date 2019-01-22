package payload

import (
	"bytes"
	"encoding/binary"
	"strings"

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
		tags := tags(fieldType.Tag)
		switch field.Kind() {
		case reflect.Ptr:
			d.pointer(field)
		case reflect.String:
			d.string(value, field, tags)
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			d.uint(value, field, tags, &bitmaskBytes)
		case reflect.Array:
			d.array(field, tags)
		case reflect.Slice:
			d.slice(value, field, tags)
		}
	}
}

func (d *decoder) slice(parent reflect.Value, value reflect.Value, tags tags) {
	length := d.dynamicLength(tags)
	value.Set(reflect.MakeSlice(value.Type(), length, length))
	for i := 0; i < length; i++ {
		sliceElement := value.Index(i)
		switch sliceElement.Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			d.uint(parent, sliceElement, tags, nil)
		case reflect.String:
			d.string(parent, sliceElement, tags)
		case reflect.Ptr:
			d.pointer(sliceElement)
		case reflect.Struct:
			d.strukt(sliceElement)
		}
	}
}

func (d *decoder) array(value reflect.Value, tags tags) {
	if value.Len() > 0 {
		size := int(value.Index(0).Type().Size())
		for i := 0; i < value.Len(); i++ {
			arrayElem := value.Index(i)
			v := d.readUint(tags.endianness(), size)
			arrayElem.SetUint(v)
		}
	}
}

func (d *decoder) uint(parent reflect.Value, value reflect.Value, tags tags, bitmaskBytes *uint64) {
	if tags.cond().nonEmpty() {
		if !checkCondition(tags.cond(), parent) {
			return
		}
	}
	if value.CanAddr() {
		if tags.bits().nonEmpty() {
			if tags.bitmask() == "start" {
				*bitmaskBytes = d.readUint(tags.endianness(), int(value.Type().Size()))
			}
			bitmaskBits := bitmaskBits(tags.bits())
			pos := util.FirstBitPosition(bitmaskBits)
			v := (*bitmaskBytes & bitmaskBits) >> pos
			value.SetUint(v)
		} else if tags.bound().nonEmpty() {
			size, _ := strconv.Atoi(string(tags.bound()))
			v := d.readUint(tags.endianness(), size)
			value.SetUint(v)
		} else {
			v := d.readUint(tags.endianness(), int(value.Type().Size()))
			value.SetUint(v)
		}
	} else {
		panic("Unaddressable uint value")
	}
}

func checkCondition(cond tag, parent reflect.Value) bool {
	v := strings.Split(string(cond), ":")
	t := v[0]
	c := v[1]
	var op string
	switch {
	case strings.Contains(c, "=="):
		op = "=="
	case strings.Contains(c, "!="):
		op = "!="
	}
	v = strings.Split(c, op)
	l := v[0]
	r := v[1]
	switch t {
	case "uint":
		lv := uint64(parent.FieldByName(l).Uint())
		n, _ := strconv.Atoi(r)
		rv := uint64(n)
		switch op {
		case "==":
			return lv == rv
		case "!=":
			return lv != rv
		}
	}
	return true
}

func (d *decoder) string(parent reflect.Value, value reflect.Value, tags tags) {
	if tags.cond().nonEmpty() {
		if !checkCondition(tags.cond(), parent) {
			return
		}
	}
	if tags.hex().nonEmpty() {
		size, _ := strconv.Atoi(string(tags.hex()))
		v := d.readUint(tags.endianness(), size)
		hexString, _ := util.UintToHexString(v, size)
		value.SetString(hexString)
	} else {
		length := d.dynamicLength(tags)
		b := make([]uint8, length, length)
		d.read(tags.endianness(), b)
		value.SetString(string(b))
	}
}

func (d *decoder) dynamicLength(tags tags) int {
	if tags.size().nonEmpty() {
		size, _ := strconv.Atoi(string(tags.size()))
		return int(d.readUint(tags.endianness(), size))
	}
	return len(d.buf.Bytes())
}

func (d *decoder) read(endianness tag, v interface{}) {
	binary.Read(d.buf, order(endianness), v)
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
