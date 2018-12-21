package pkg

import (
	"github.com/petomalina/mirror"
	"reflect"
	"unsafe"
)

type HijackedField struct {
}

func Hijack(model *mirror.Struct) error {

	for _, f := range model.RawFields() {
		// ignore exported fields - user already has access to them
		if f.Exported() {
			continue
		}

		field := reflect.NewAt(f.Typ, unsafe.Pointer(f.Value.UnsafeAddr())).Elem()
		field.Set(reflect.ValueOf("Hello"))
	}

	return nil
}
