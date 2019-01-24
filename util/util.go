package util

import (
	"fmt"
)

func Panicf(format string, v ...interface{}) {
	panic(fmt.Sprintf(format, v...))
}
