package mirror

import (
	"reflect"
)

type Struct struct {
	Ref reflect.Value
}

func ReflectStruct(s interface{}) *Struct {
	return &Struct{reflect.ValueOf(s)}
}

func ReflectStructs(ss ...interface{}) []*Struct {
	rss := []*Struct{}

	for _, s := range ss {
		rss = append(rss, ReflectStruct(s))
	}

	return rss
}

func (s *Struct) Name() string {
	return s.Ref.Type().Name()
}

// Fields returns a map of Name:Type pairs which can be used
// directly when generating new code
func (s *Struct) Fields() map[string]string {
	ff := map[string]string{}

	for _, f := range s.RawFields() {
		ff[f.Field.Name] = f.Typ.String()
	}

	return ff
}

func (s *Struct) IsInterface(i reflect.Type) bool {
	return false
}

// RawStructFieldType is returned by the RawFields as an array of
// key-value pairs (as map can't be indexed by full structure and
// reflect.StuctField is used without pointers across reflect)
type RawStructFieldType struct {
	Field reflect.StructField
	Typ   reflect.Type
}

func (s *Struct) RawFields() []RawStructFieldType {
	rf := []RawStructFieldType{}

	num := s.Ref.NumField()
	for i := 0; i < num; i++ {
		v := s.Ref.Field(i)

		rf = append(rf, RawStructFieldType{
			Field: s.Ref.Type().Field(i),
			Typ: v.Type(),
		})
	}

	return rf
}
