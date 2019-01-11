package reflection

import (
	"reflect"

	"github.com/dyrkin/znp-go/util"
)

func Copy(n interface{}) interface{} {
	v := reflect.ValueOf(n)
	switch v.Kind() {
	case reflect.Struct:
		copy := reflect.New(v.Type()).Elem()
		return copy.Interface()
	case reflect.Ptr:
		e := v.Elem()
		copy := reflect.New(e.Type())
		return copy.Interface()
	}
	util.Panicf("reflection.Copy: Unsupported value: %#v", n)
	return nil
}
