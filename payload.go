package znp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
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
	buf := &bytes.Buffer{}
	mirror := reflect.ValueOf(request).Elem()
	for i := 0; i < mirror.NumField(); i++ {
		valueMirror := mirror.Field(i)
		typeMirror := mirror.Type().Field(i)
		tagMirror := typeMirror.Tag
		subtype := tagMirror.Get("subtype")
		endianness := tagMirror.Get("endianness")
		switch value := valueMirror.Interface().(type) {
		case string:
			switch subtype {
			case "longaddr":
				addr, _ := strconv.ParseInt(value[2:], 16, 64)
				write(buf, endianness, addr)
			}
		case []*Network:
			for _, v := range value {
				writeNetwork(buf, v)
			}
		default:
			write(buf, endianness, value)
		}
	}
	return buf.Bytes()
}

func deserialize(payload []byte, response interface{}) {
	if len(payload) == 0 {
		return
	}
	buf := bytes.NewBuffer(payload)
	mirror := reflect.ValueOf(response).Elem()
	var dynBufLen uint64
	for i := 0; i < mirror.NumField(); i++ {
		valueMirror := mirror.Field(i)
		typeMirror := mirror.Type().Field(i)
		tagMirror := typeMirror.Tag
		subtype := tagMirror.Get("subtype")
		endianness := tagMirror.Get("endianness")
		switch value := valueMirror.Interface().(type) {
		case string:
			switch subtype {
			case "longaddr":
				var v uint64
				read(buf, endianness, &v)
				valueMirror.SetString(fmt.Sprintf("0x%016x", v))
			}
		case uint8:
			var v uint8
			read(buf, endianness, &v)
			valueMirror.Set(reflect.ValueOf(v))
			dynBufLen = uint64(v)
		case uint16:
			var v uint16
			read(buf, endianness, &v)
			valueMirror.Set(reflect.ValueOf(v))
			dynBufLen = uint64(v)
		case uint32:
			var v uint32
			read(buf, endianness, &v)
			valueMirror.Set(reflect.ValueOf(v))
			dynBufLen = uint64(v)
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
		case []*Network:
			v := make([]*Network, dynBufLen)
			for i := range v {
				v[i] = readNetwork(buf)
			}
			valueMirror.Set(reflect.ValueOf(v))
		default:
			log.Printf("Unsupported type: %+v", value)
		}
	}
}
