package znp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"
)

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

func serialize(request interface{}) []byte {
	if request == nil {
		return make([]byte, 0)
	}
	bitmaskStarted := false
	bitmaskStopped := false
	var bitmaskBytes uint64
	buf := &bytes.Buffer{}
	mirror := reflect.ValueOf(request).Elem()
	for i := 0; i < mirror.NumField(); i++ {
		valueMirror := mirror.Field(i)
		typeMirror := mirror.Type().Field(i)
		tagMirror := typeMirror.Tag
		hex := tagMirror.Get("hex")
		endianness := tagMirror.Get("endianness")
		bitmask := tagMirror.Get("bitmask")
		if bitmask != "" {
			bitmaskStarted = bitmask == "start"
			bitmaskStopped = !bitmaskStarted
			if bitmaskStarted {
				bitmaskBytes = 0
			}
		}
		var processBitmask = func(v interface{}) bool {
			if bitmaskStarted || bitmaskStopped {
				bits := tagMirror.Get("bits")
				if bits == "" && bitmaskStarted {
					log.Fatalf("Bitmask is started but bits tag is not defined")
				}
				var bitmaskBits uint64
				if strings.HasPrefix(bits, "0x") {
					bitmaskBits, _ = strconv.ParseUint(bits[2:], 16, len(bits[2:])*4)
				} else if strings.HasPrefix(bits, "0b") {
					bitmaskBits, _ = strconv.ParseUint(bits[2:], 2, len(bits[2:]))
				}
				pos := getFirstSetBitPos(bitmaskBits)
				bitmaskBytes = bitmaskBytes | ((toUint64(v) << pos) & bitmaskBits)
				if !bitmaskStarted && bitmaskStopped {
					write(buf, endianness, toVtype(v, bitmaskBytes))
				}

				bitmaskStopped = false
				return true
			}
			return false
		}
		len := tagMirror.Get("len")
		if len != "" {
			switch len {
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
		case uint8:
			if !processBitmask(value) {
				write(buf, endianness, value)
			}
		case uint16:
			if !processBitmask(value) {
				write(buf, endianness, value)
			}
		case uint32:
			if !processBitmask(value) {
				write(buf, endianness, value)
			}
		case uint64:
			if !processBitmask(value) {
				write(buf, endianness, value)
			}
		case string:
			switch hex {
			case "uint64":
				addr, _ := strconv.ParseUint(value[2:], 16, 64)
				write(buf, endianness, addr)
			case "uint32":
				addr, _ := strconv.ParseUint(value[2:], 16, 32)
				write(buf, endianness, uint32(addr))
			case "uint16":
				addr, _ := strconv.ParseUint(value[2:], 16, 16)
				write(buf, endianness, uint16(addr))
			case "uint8":
				addr, _ := strconv.ParseUint(value[2:], 16, 8)
				write(buf, endianness, uint8(addr))
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
	bitmaskStarted := false
	bitmaskStopped := false
	var bitmaskBytes uint64
	for i := 0; i < mirror.NumField(); i++ {
		valueMirror := mirror.Field(i)
		typeMirror := mirror.Type().Field(i)
		tagMirror := typeMirror.Tag
		hex := tagMirror.Get("hex")
		endianness := tagMirror.Get("endianness")
		length := tagMirror.Get("len")
		bitmask := tagMirror.Get("bitmask")
		bitmaskStartedNow := false
		if bitmask != "" {
			bitmaskStarted = bitmask == "start"
			bitmaskStopped = !bitmaskStarted
			bitmaskStartedNow = bitmaskStarted
		}
		var processBitmask = func(v interface{}) bool {
			if bitmaskStartedNow {
				read(buf, endianness, v)
				bitmaskBytes = toUint64(v)
			}
			if bitmaskStarted || bitmaskStopped {
				bits := tagMirror.Get("bits")
				if bits == "" && bitmaskStarted {
					log.Fatalf("Bitmask is started but bits tag is not defined")
				}
				var bitmaskBits uint64
				if strings.HasPrefix(bits, "0x") {
					bitmaskBits, _ = strconv.ParseUint(bits[2:], 16, len(bits[2:])*4)
				} else if strings.HasPrefix(bits, "0b") {
					bitmaskBits, _ = strconv.ParseUint(bits[2:], 2, len(bits[2:]))
				}
				pos := getFirstSetBitPos(bitmaskBits)
				v := toVtype(v, bitmaskBytes&bitmaskBits>>pos)
				valueMirror.Set(reflect.ValueOf(v))
				bitmaskStopped = false
				return true
			}
			return false
		}
		var dynBufLen uint32
		if length != "" {
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
		switch valueMirror.Interface().(type) {
		case string:
			switch hex {
			case "uint64":
				var v uint64
				read(buf, endianness, &v)
				valueMirror.SetString(fmt.Sprintf("0x%016x", v))
			case "uint32":
				var v uint32
				read(buf, endianness, &v)
				valueMirror.SetString(fmt.Sprintf("0x%08x", v))
			case "uint16":
				var v uint16
				read(buf, endianness, &v)
				valueMirror.SetString(fmt.Sprintf("0x%04x", v))
			case "uint8":
				var v uint8
				read(buf, endianness, &v)
				valueMirror.SetString(fmt.Sprintf("0x%02x", v))
			}
		case uint8:
			var v uint8
			if !processBitmask(&v) {
				read(buf, endianness, &v)
				valueMirror.Set(reflect.ValueOf(v))
			}
		case uint16:
			var v uint16
			if !processBitmask(&v) {
				read(buf, endianness, &v)
				valueMirror.Set(reflect.ValueOf(v))
			}
		case uint32:
			var v uint32
			if !processBitmask(&v) {
				read(buf, endianness, &v)
				valueMirror.Set(reflect.ValueOf(v))
			}
		case uint64:
			var v uint64
			if !processBitmask(&v) {
				read(buf, endianness, &v)
				valueMirror.Set(reflect.ValueOf(v))
			}
		case [8]byte:
			var v [8]byte
			read(buf, endianness, &v)
			valueMirror.Set(reflect.ValueOf(v))
		case [16]byte:
			var v [16]byte
			read(buf, endianness, &v)
			valueMirror.Set(reflect.ValueOf(v))
		case [18]byte:
			var v [18]byte
			read(buf, endianness, &v)
			valueMirror.Set(reflect.ValueOf(v))
		case [32]byte:
			var v [32]byte
			read(buf, endianness, &v)
			valueMirror.Set(reflect.ValueOf(v))
		case [42]byte:
			var v [42]byte
			read(buf, endianness, &v)
			valueMirror.Set(reflect.ValueOf(v))
		case [100]byte:
			var v [100]byte
			read(buf, endianness, &v)
			valueMirror.Set(reflect.ValueOf(v))
		case [2]uint16:
			var v [2]uint16
			read(buf, endianness, &v)
			valueMirror.Set(reflect.ValueOf(v))
		case []byte:
			v := make([]byte, dynBufLen)
			read(buf, endianness, v)
			valueMirror.Set(reflect.ValueOf(v))
		case []uint16:
			v := make([]uint16, dynBufLen)
			read(buf, endianness, v)
			valueMirror.Set(reflect.ValueOf(v))
		default:
			switch valueMirror.Kind() {
			case reflect.Ptr:
				el := reflect.New(valueMirror.Type().Elem())
				v := el.Interface()
				if valueMirror.CanSet() {
					valueMirror.Set(el)
				}
				deserialize(buf, v)
			case reflect.Slice:
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

func getFirstSetBitPos(n uint64) uint8 {
	return uint8(math.Log2(float64(n & -n)))
}

func toUint64(v interface{}) uint64 {
	switch z := v.(type) {
	case *uint8:
		return uint64(*z)
	case *uint16:
		return uint64(*z)
	case *uint32:
		return uint64(*z)
	case *uint64:
		return uint64(*z)
	case uint8:
		return uint64(z)
	case uint16:
		return uint64(z)
	case uint32:
		return uint64(z)
	case uint64:
		return uint64(z)
	}
	return 0
}

func toVtype(v interface{}, val uint64) interface{} {
	switch v.(type) {
	case *uint8, uint8:
		return uint8(val)
	case *uint16, uint16:
		return uint16(val)
	case *uint32, uint32:
		return uint32(val)
	case *uint64, uint64:
		return uint64(val)
	}
	log.Fatalf("Unknown value: %v", v)
	return 0
}
