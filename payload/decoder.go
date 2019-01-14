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
	for i := 0; i < value.Len(); i++ {
		arrayElement := value.Index(i)
		arrayElementNew := reflect.New(arrayElement.Type())
		v := arrayElementNew.Interface()
		d.read(tags.endianness, v)
		arrayElement.Set(reflect.ValueOf(v).Elem())
	}
}

func (d *decoder) uint(value reflect.Value, tags *tags, bitmaskBytes *uint64) {
	if value.CanAddr() {
		ptr := value.Addr()
		if tags.bits.nonEmpty() {
			if tags.bitmask == "start" {
				d.read(tags.endianness, ptr.Interface())
				*bitmaskBytes = valueConvertTo(value, uint64Type).Interface().(uint64)
			}
			bitmaskBits := bitmaskBits(tags.bits)
			pos := util.FirstBitPosition(bitmaskBits)
			v := (*bitmaskBytes & bitmaskBits) >> pos
			value.Set(valueConvertTo(reflect.ValueOf(v), value.Type()))
		} else if tags.bound.nonEmpty() {
			size, _ := strconv.Atoi(string(tags.bound))
			a := make([]uint8, size, size)
			endianness := tags.endianness
			d.read(endianness, a)
			v := uint64(0)
			for i := 0; i < size; i++ {
				if endianness == "be" {
					v = v | (uint64(a[size-i-1]) << byte(i*8))
				} else {
					v = v | (uint64(a[i]) << byte(i*8))
				}
			}
			value.Set(valueConvertTo(reflect.ValueOf(v), value.Type()))
		} else {
			v := ptr.Interface()
			d.read(tags.endianness, v)
		}
	} else {
		panic("Unaddressable uint value")
	}
}

func (d *decoder) string(value reflect.Value, tags *tags) {
	if tags.hex.nonEmpty() {
		size, _ := strconv.Atoi(string(tags.hex))
		if typ, ok := types[size]; ok {
			v := reflect.New(typ)
			ptr := v.Interface()
			d.read(tags.endianness, ptr)
			hexString, _ := util.UintToHexString(v.Elem().Interface())
			value.SetString(hexString)
		} else {
			util.Panicf("Unsupported hex size: %s", tags.hex)
		}
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
		if typ, ok := types[size]; ok {
			v := reflect.New(typ)
			ptr := v.Interface()
			d.read(tags.endianness, ptr)
			return valueConvertTo(v.Elem(), intType).Interface().(int)
		}
		util.Panicf("Unsupported length: %s", tags.size)
		return 0
	} else {
		return len(d.buf.Bytes())
	}
}

func (d *decoder) read(endianness tag, v interface{}) {
	binary.Read(d.buf, order(endianness), v)
}
