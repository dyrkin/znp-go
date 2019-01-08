package util

import (
	"log"
	"math"
)

func Uint64(v interface{}) uint64 {
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
	default:
		log.Fatalf("Unsupported value: %#v: ", v)
		return 0
	}
}

func Vtype(v interface{}, val uint64) interface{} {
	switch v.(type) {
	case *uint8, uint8:
		return uint8(val)
	case *uint16, uint16:
		return uint16(val)
	case *uint32, uint32:
		return uint32(val)
	case *uint64, uint64:
		return uint64(val)
	default:
		log.Fatalf("Unsupported value: %#v", v)
		return 0
	}
}

func FirstBitPosition(n uint64) uint8 {
	return uint8(math.Log2(float64(n & -n)))
}
