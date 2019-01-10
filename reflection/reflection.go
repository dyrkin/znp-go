package reflection

import (
	"log"
	"reflect"
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
	log.Fatalf("Unsupported value: %#v", n)
	return nil
}

func Init(r interface{}) {
	m := reflect.ValueOf(r).Elem()
	el := reflect.New(m.Type().Elem())
	if m.CanSet() {
		m.Set(el)
	}
}

func GetTag(tags reflect.StructTag, tagName string) (string, bool) {
	v := tags.Get(tagName)
	return v, v != ""
}
