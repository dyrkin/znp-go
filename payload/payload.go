package payload

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/dyrkin/znp-go/util"
)

var types = map[int]reflect.Type{
	1: reflect.TypeOf(uint8(0)),
	2: reflect.TypeOf(uint16(0)),
	4: reflect.TypeOf(uint32(0)),
	8: reflect.TypeOf(uint64(0)),
}

var intType = reflect.TypeOf(int(0))
var uint64Type = reflect.TypeOf(uint64(0))

func Encode(request interface{}) []byte {
	value := reflect.ValueOf(request)
	buf := &bytes.Buffer{}
	encode(buf, value)
	return buf.Bytes()
}

func Decode(payload []byte, response interface{}) {
	value := reflect.ValueOf(response)
	decode(bytes.NewBuffer(payload), value)
}

func encode(buf *bytes.Buffer, value reflect.Value) {
	switch value.Kind() {
	case reflect.Ptr:
		encodePointer(buf, value)
	case reflect.Struct:
		encodeStruct(buf, value)
	}
}

func encodeStruct(buf *bytes.Buffer, value reflect.Value) {
	var bitmaskBytes uint64
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := value.Type().Field(i)
		tags := newTags(fieldType)
		switch field.Kind() {
		case reflect.Ptr:
			encodePointer(buf, field)
		case reflect.String:
			encodeString(buf, field, tags)
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			encodeUint(buf, field, tags, &bitmaskBytes)
		case reflect.Array:
			encodeArray(buf, field, tags)
		case reflect.Slice:
			encodeSlice(buf, field, tags)
		case reflect.Interface:
			encodePointer(buf, field)
		}
	}
}

func encodeSlice(buf *bytes.Buffer, value reflect.Value, tags *tags) {
	length := value.Len()
	writeDynamicLength(buf, length, tags)
	for i := 0; i < length; i++ {
		sliceElement := value.Index(i)
		switch sliceElement.Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			encodeUint(buf, sliceElement, tags, nil)
		case reflect.String:
			encodeString(buf, sliceElement, tags)
		case reflect.Ptr:
			encodePointer(buf, sliceElement)
		case reflect.Struct:
			encodeStruct(buf, sliceElement)
		}
	}
}

func encodeArray(buf *bytes.Buffer, value reflect.Value, tags *tags) {
	write(buf, tags.endianness.value, value.Interface())
}

func encodeString(buf *bytes.Buffer, value reflect.Value, tags *tags) {
	v := value.Interface().(string)
	if tags.hex.nonEmpty() {
		size, _ := strconv.Atoi(tags.hex.value)
		typ := types[size]
		addr, _ := strconv.ParseUint(v[2:], 16, size*8)
		v := convertTo(addr, typ)
		write(buf, tags.endianness.value, v)
	} else {
		writeDynamicLength(buf, len(v), tags)
		write(buf, tags.endianness.value, []uint8(v))
	}
}

func encodeUint(buf *bytes.Buffer, value reflect.Value, tags *tags, bitmaskBytes *uint64) {
	if tags.bits.nonEmpty() {
		if tags.bitmask.value == "start" {
			*bitmaskBytes = 0
		}
		bitmaskBits := bitmaskBits(tags.bits.value)
		pos := util.FirstBitPosition(bitmaskBits)
		v := valueConvertTo(value, uint64Type).Interface().(uint64)
		v = ((v << pos) & bitmaskBits)
		*bitmaskBytes = (*bitmaskBytes) | v
		if tags.bitmask.value == "end" {
			v := convertTo(*bitmaskBytes, value.Type())
			write(buf, tags.endianness.value, v)
		}

	} else if tags.bound.nonEmpty() {
		size, _ := strconv.Atoi(tags.bound.value)
		v := make([]uint8, size, size)
		uintVal := value.Uint()
		endianness := tags.endianness.value
		for i := 0; i < size; i++ {
			if endianness == "be" {
				v[size-i-1] = byte(uintVal >> byte(i*8))
			} else {
				v[i] = byte(uintVal >> byte(i*8))
			}

		}
		write(buf, tags.endianness.value, v)
	} else {
		v := value.Interface()
		write(buf, tags.endianness.value, v)
	}
}

func encodePointer(buf *bytes.Buffer, value reflect.Value) {
	encode(buf, value.Elem())
}

func decode(buf *bytes.Buffer, value reflect.Value) {
	switch value.Kind() {
	case reflect.Ptr:
		decodePointer(buf, value)
	case reflect.Struct:
		decodeStruct(buf, value)
	}
}

func decodePointer(buf *bytes.Buffer, value reflect.Value) {
	if value.IsNil() {
		element := reflect.New(value.Type().Elem())
		if value.CanSet() {
			value.Set(element)
		}
	}
	decode(buf, value.Elem())
}

func decodeStruct(buf *bytes.Buffer, value reflect.Value) {
	var bitmaskBytes uint64
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := value.Type().Field(i)
		tags := newTags(fieldType)
		switch field.Kind() {
		case reflect.Ptr:
			decodePointer(buf, field)
		case reflect.String:
			decodeString(buf, field, tags)
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			decodeUint(buf, field, tags, &bitmaskBytes)
		case reflect.Array:
			decodeArray(buf, field, tags)
		case reflect.Slice:
			decodeSlice(buf, field, tags)
		}
	}
}

func decodeSlice(buf *bytes.Buffer, value reflect.Value, tags *tags) {
	length := readDynamicLength(buf, tags)
	value.Set(reflect.MakeSlice(value.Type(), length, length))
	for i := 0; i < length; i++ {
		sliceElement := value.Index(i)
		switch sliceElement.Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			decodeUint(buf, sliceElement, tags, nil)
		case reflect.String:
			decodeString(buf, sliceElement, tags)
		case reflect.Ptr:
			decodePointer(buf, sliceElement)
		case reflect.Struct:
			decodeStruct(buf, sliceElement)
		}
	}
}

func decodeArray(buf *bytes.Buffer, value reflect.Value, tags *tags) {
	for i := 0; i < value.Len(); i++ {
		arrayElement := value.Index(i)
		arrayElementNew := reflect.New(arrayElement.Type())
		v := arrayElementNew.Interface()
		read(buf, tags.endianness.value, v)
		arrayElement.Set(reflect.ValueOf(v).Elem())
	}
}

func decodeUint(buf *bytes.Buffer, value reflect.Value, tags *tags, bitmaskBytes *uint64) {
	if value.CanAddr() {
		ptr := value.Addr()
		if tags.bits.nonEmpty() {
			if tags.bitmask.value == "start" {
				read(buf, tags.endianness.value, ptr.Interface())
				*bitmaskBytes = valueConvertTo(value, uint64Type).Interface().(uint64)
			}
			bitmaskBits := bitmaskBits(tags.bits.value)
			pos := util.FirstBitPosition(bitmaskBits)
			v := (*bitmaskBytes & bitmaskBits) >> pos
			value.Set(valueConvertTo(reflect.ValueOf(v), value.Type()))
		} else if tags.bound.nonEmpty() {
			size, _ := strconv.Atoi(tags.bound.value)
			a := make([]uint8, size, size)
			endianness := tags.endianness.value
			read(buf, endianness, a)
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
			read(buf, tags.endianness.value, v)
		}
	} else {
		panic("Unaddressable uint value")
	}
}

func decodeString(buf *bytes.Buffer, value reflect.Value, tags *tags) {
	if tags.hex.nonEmpty() {
		size, _ := strconv.Atoi(tags.hex.value)
		if typ, ok := types[size]; ok {
			v := reflect.New(typ)
			ptr := v.Interface()
			read(buf, tags.endianness.value, ptr)
			hexString, _ := util.UintToHexString(v.Elem().Interface())
			value.SetString(hexString)
		} else {
			util.Panicf("Unsupported hex size: %s", tags.hex.value)
		}
	} else {
		length := readDynamicLength(buf, tags)
		b := make([]uint8, length, length)
		read(buf, tags.endianness.value, b)
		value.SetString(string(b))
	}
}

func readDynamicLength(buf *bytes.Buffer, tags *tags) int {
	if tags.size.nonEmpty() {
		size, _ := strconv.Atoi(tags.size.value)
		if typ, ok := types[size]; ok {
			v := reflect.New(typ)
			ptr := v.Interface()
			read(buf, tags.endianness.value, ptr)
			return convertTo(v.Elem().Interface(), intType).(int)
		}
		util.Panicf("Unsupported length: %s", tags.size.value)
		return 0
	} else {
		return len(buf.Bytes())
	}
}

func writeDynamicLength(buf *bytes.Buffer, length int, tags *tags) {
	if tags.size.nonEmpty() {
		size, _ := strconv.Atoi(tags.size.value)
		if typ, ok := types[size]; ok {
			v := reflect.New(typ).Elem()
			v.Set(valueConvertTo(reflect.ValueOf(length), typ))
			write(buf, tags.endianness.value, v.Interface())
		}
	}
}

type tag struct {
	value string
}

func (t *tag) nonEmpty() bool {
	return len(t.value) > 0
}

type tags struct {
	hex        *tag
	endianness *tag
	size       *tag
	bitmask    *tag
	bits       *tag
	bound      *tag
}

func newTags(field reflect.StructField) *tags {
	hex := &tag{field.Tag.Get("hex")}
	endianness := &tag{field.Tag.Get("endianness")}
	size := &tag{field.Tag.Get("size")}
	bitmask := &tag{field.Tag.Get("bitmask")}
	bits := &tag{field.Tag.Get("bits")}
	bound := &tag{field.Tag.Get("bound")}
	return &tags{hex: hex,
		endianness: endianness,
		size:       size,
		bitmask:    bitmask,
		bits:       bits,
		bound:      bound,
	}
}

func convertTo(v interface{}, typ reflect.Type) interface{} {
	value := reflect.ValueOf(v)
	return valueConvertTo(value, typ).Interface()
}

func valueConvertTo(value reflect.Value, typ reflect.Type) reflect.Value {
	return value.Convert(typ)
}

func bitmaskBits(value string) uint64 {
	var bitmaskBits uint64
	if strings.HasPrefix(value, "0x") {
		bitmaskBits, _ = strconv.ParseUint(value[2:], 16, len(value[2:])*4)
	} else if strings.HasPrefix(value, "0b") {
		bitmaskBits, _ = strconv.ParseUint(value[2:], 2, len(value[2:]))
	}
	return bitmaskBits
}

func order(endianness string) binary.ByteOrder {
	if endianness == "be" {
		return binary.BigEndian
	}
	return binary.LittleEndian
}

func write(w io.Writer, endianness string, v interface{}) {
	binary.Write(w, order(endianness), v)
}

func read(r io.Reader, endianness string, v interface{}) {
	binary.Read(r, order(endianness), v)
}
