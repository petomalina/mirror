package pkg

import (
	"fmt"
	"github.com/petomalina/mirror"
	"reflect"
	"unsafe"
)

type HijackedField struct {
}

func Hijack(model interface{}) error {
	fields := mirror.ReflectStruct(model).RawFields()

	fmt.Printf("%+v\n", model)

	for _, f := range fields {
		// ignore exported fields - user already has access to them
		if f.Exported() {
			continue
		}

		field := reflect.NewAt(f.Typ, unsafe.Pointer(f.Value.UnsafeAddr())).Elem()
		field.Set(reflect.ValueOf("Hello"))
	}

	return nil
}
