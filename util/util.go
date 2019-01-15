package util

import (
	"fmt"
	"math"
)

func FirstBitPosition(n uint64) uint8 {
	return uint8(math.Log2(float64(n & -n)))
}

func UintToHexString(v uint64, size int) (string, error) {
	switch size {
	case 1:
		return fmt.Sprintf("0x%02x", v), nil
	case 2:
		return fmt.Sprintf("0x%04x", v), nil
	case 4:
		return fmt.Sprintf("0x%08x", v), nil
	case 8:
		return fmt.Sprintf("0x%016x", v), nil
	}
	return "", fmt.Errorf("reflection.Vtype: Unsupported value: %s", v)
}

func Panicf(format string, v ...interface{}) {
	panic(fmt.Sprintf(format, v...))
}
