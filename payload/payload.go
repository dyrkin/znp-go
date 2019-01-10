package payload

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/dyrkin/znp-go/reflection"

	"github.com/dyrkin/znp-go/util"
)

func Encode(request interface{}) []byte {
	return serialize(request)
}

func Decode(payload []byte, response interface{}) {
	deserialize(bytes.NewBuffer(payload), response)
}

func serialize(request interface{}) []byte {
	if request == nil {
		return make([]byte, 0)
	}
	var bitmaskBytes uint64
	buf := &bytes.Buffer{}
	mirror := reflect.ValueOf(request).Elem()
	for i := 0; i < mirror.NumField(); i++ {
		valueMirror := mirror.Field(i)
		typeMirror := mirror.Type().Field(i)
		tags := typeMirror.Tag
		hex, hexOk := reflection.GetTag(tags, "hex")
		endianness, _ := reflection.GetTag(tags, "endianness")
		length, lengthOk := reflection.GetTag(tags, "len")
		bitmask, _ := reflection.GetTag(tags, "bitmask")
		bits, bitsOk := reflection.GetTag(tags, "bits")
		var writeBitmask = func(v interface{}) bool {
			if bitsOk {
				if bitmask == "start" {
					bitmaskBytes = 0
				}
				bitmaskBits := bitmaskBits(bits)
				pos := util.FirstBitPosition(bitmaskBits)
				bitmaskBytes = bitmaskBytes | ((util.Uint64(v) << pos) & bitmaskBits)
				if bitmask == "end" {
					write(buf, endianness, util.Vtype(v, bitmaskBytes))
				}
				return true
			}
			return false
		}
		var writeString = func(v string) {
			if hexOk {
				switch hex {
				case "uint64":
					addr, _ := strconv.ParseUint(v[2:], 16, 64)
					write(buf, endianness, addr)
				case "uint32":
					addr, _ := strconv.ParseUint(v[2:], 16, 32)
					write(buf, endianness, uint32(addr))
				case "uint16":
					addr, _ := strconv.ParseUint(v[2:], 16, 16)
					write(buf, endianness, uint16(addr))
				case "uint8":
					addr, _ := strconv.ParseUint(v[2:], 16, 8)
					write(buf, endianness, uint8(addr))
				}
			} else {
				write(buf, endianness, []uint8(v))
			}
		}
		writeUint := func(v interface{}) {
			if !writeBitmask(v) {
				write(buf, endianness, v)
			}
		}
		if lengthOk {
			switch length {
			case "uint8":
				v := uint8(valueMirror.Len())
				write(buf, endianness, v)
			case "uint16":
				v := uint16(valueMirror.Len())
				write(buf, endianness, v)
			case "uint32":
				v := uint32(valueMirror.Len())
				write(buf, endianness, v)
			}
		}
		switch value := valueMirror.Interface().(type) {
		case uint8, uint16, uint32, uint64:
			writeUint(value)
		case string:
			writeString(value)
		case []string:
			for _, v := range value {
				writeString(v)
			}
		default:
			switch {
			case valueMirror.Kind() == reflect.Ptr:
				write(buf, endianness, serialize(value))
			case valueMirror.Kind() == reflect.Slice && valueMirror.Len() > 0 && valueMirror.Index(0).Kind() == reflect.Ptr:
				for i := 0; i < valueMirror.Len(); i++ {
					write(buf, endianness, serialize(valueMirror.Index(i).Interface()))
				}
			default:
				write(buf, endianness, value)
			}
		}
	}
	return buf.Bytes()
}

func deserialize(buf *bytes.Buffer, response interface{}) {
	mirror := reflect.ValueOf(response).Elem()
	if (mirror.Kind() == reflect.Ptr) && mirror.IsNil() {
		reflection.Init(response)
		mirror = reflect.ValueOf(response).Elem().Elem()
	}
	var bitmaskBytes uint64
	for i := 0; i < mirror.NumField(); i++ {
		valueMirror := mirror.Field(i)
		typeMirror := mirror.Type().Field(i)
		tags := typeMirror.Tag
		hex, hexOk := reflection.GetTag(tags, "hex")
		endianness, _ := reflection.GetTag(tags, "endianness")
		length, lengthOk := reflection.GetTag(tags, "len")
		bitmask, _ := reflection.GetTag(tags, "bitmask")
		bits, bitsOk := reflection.GetTag(tags, "bits")
		var setBitmask = func(v interface{}) bool {
			if bitsOk {
				if bitmask == "start" {
					read(buf, endianness, v)
					bitmaskBytes = util.Uint64(v)
				}
				bitmaskBits := bitmaskBits(bits)
				pos := util.FirstBitPosition(bitmaskBits)
				v := util.Vtype(v, bitmaskBytes&bitmaskBits>>pos)
				valueMirror.Set(reflect.ValueOf(v).Convert(valueMirror.Type()))
				return true
			}
			return false
		}
		readHexString := func() string {
			readHex := func(v interface{}) string {
				read(buf, endianness, v)
				hexString, _ := util.UintToHexString(v)
				return hexString
			}
			switch hex {
			case "uint64":
				var v uint64
				return readHex(&v)
			case "uint32":
				var v uint32
				return readHex(&v)
			case "uint16":
				var v uint16
				return readHex(&v)
			case "uint8":
				var v uint8
				return readHex(&v)
			default:
				log.Fatalf("Unsupported hex size: %s", hex)
			}
			return ""
		}
		var dynBufLen uint32
		if lengthOk {
			switch length {
			case "uint8":
				var v uint8
				read(buf, endianness, &v)
				dynBufLen = uint32(v)
			case "uint16":
				var v uint16
				read(buf, endianness, &v)
				dynBufLen = uint32(v)
			case "uint32":
				var v uint32
				read(buf, endianness, &v)
				dynBufLen = v
			}
		}
		setUint := func(v interface{}) {
			if !setBitmask(v) {
				read(buf, endianness, v)
				valueMirror.Set(reflect.ValueOf(v).Elem().Convert(valueMirror.Type()))
			}
		}
		setSlice := func() {
			if valueMirror.CanSet() {
				valueMirror.Set(reflect.MakeSlice(valueMirror.Type(), int(dynBufLen), int(dynBufLen)))
				for i := 0; i < int(dynBufLen); i++ {
					sliceElemMirror := valueMirror.Index(i)
					el := reflect.New(sliceElemMirror.Type().Elem())
					v := el.Interface()
					deserialize(buf, v)
					sliceElemMirror.Set(reflect.ValueOf(v))
				}
			}
		}
		setArray := func() {
			for i := 0; i < valueMirror.Len(); i++ {
				sliceElemMirror := valueMirror.Index(i)
				el := reflect.New(sliceElemMirror.Type())
				v := el.Interface()
				read(buf, endianness, v)
				sliceElemMirror.Set(reflect.ValueOf(v).Elem())
			}
		}
		switch valueMirror.Interface().(type) {
		case string:
			if hexOk {
				valueMirror.SetString(readHexString())
			} else {
				b := make([]uint8, dynBufLen, dynBufLen)
				read(buf, endianness, b)
				valueMirror.SetString(string(b))
			}
		case []string:
			if valueMirror.CanSet() {
				valueMirror.Set(reflect.MakeSlice(valueMirror.Type(), int(dynBufLen), int(dynBufLen)))
				for i := 0; i < int(dynBufLen); i++ {
					sliceElemMirror := valueMirror.Index(i)
					hexString := readHexString()
					sliceElemMirror.SetString(hexString)
				}
			}
		case []uint8:
			if lengthOk {
				v := make([]uint8, dynBufLen)
				read(buf, endianness, v)
				valueMirror.Set(reflect.ValueOf(v))
			} else {
				valueMirror.Set(reflect.ValueOf(buf.Bytes()))
			}
		case []uint16:
			var v []uint16
			if lengthOk {
				v = make([]uint16, dynBufLen)
			} else {
				v = make([]uint16, len(buf.Bytes())/2)
			}
			read(buf, endianness, v)
			valueMirror.Set(reflect.ValueOf(v))
		default:
			switch valueMirror.Kind() {
			case reflect.Uint8:
				var v uint8
				setUint(&v)
			case reflect.Uint16:
				var v uint16
				setUint(&v)
			case reflect.Uint32:
				var v uint32
				setUint(&v)
			case reflect.Uint64:
				var v uint64
				setUint(&v)
			case reflect.Ptr:
				el := reflect.New(valueMirror.Type().Elem())
				v := el.Interface()
				if valueMirror.CanSet() {
					valueMirror.Set(el)
				}
				deserialize(buf, v)
			case reflect.Slice:
				setSlice()
			case reflect.Array:
				setArray()
			default:
				el := reflect.New(valueMirror.Type())
				v := el.Interface()
				read(buf, endianness, v)
				if valueMirror.CanSet() {
					valueMirror.Set(reflect.ValueOf(v).Elem())
				}
			}
		}
	}
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
