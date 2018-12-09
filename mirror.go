package mirror

import (
	"reflect"
	"strings"
)

type Struct struct {
	Ref interface{}
}

// ReflectStruct creates a mirror Struct that enables users
// to use mirror-enhanced reflections
func ReflectStruct(s interface{}) *Struct {
	return &Struct{
		s,
	}
}

// ReflectStructs is a plural function for ReflectStruct
func ReflectStructs(ss ...interface{}) []*Struct {
	rss := []*Struct{}

	for _, s := range ss {
		rss = append(rss, ReflectStruct(s))
	}

	return rss
}

// Name returns the name of an underlying type for a given reflection.
// Note that this does not return pointer type asterixes nor
// package names
func (s *Struct) Name() string {
	name := reflect.TypeOf(s.Ref).String()

	// strip the package prefix, as we don't want it explicitly in the name
	pkgStrip := strings.Split(name, ".")
	if len(pkgStrip) > 1 {
		return pkgStrip[1]
	}

	return pkgStrip[0]
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

// IsInterface returns true if the reflection complies with
// the given interface value
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

// Exported returns true if the field is exported within the struct
func (f *RawStructFieldType) Exported() bool {
	return strings.Title(f.Field.Name) == f.Field.Name
}

// RawFields returns all fields of a given reflection type
func (s *Struct) RawFields() []RawStructFieldType {
	rf := []RawStructFieldType{}

	sValue := reflect.ValueOf(s.Ref)

	num := sValue.NumField()
	for i := 0; i < num; i++ {
		v := sValue.Field(i)

		rf = append(rf, RawStructFieldType{
			Field: sValue.Type().Field(i),
			Typ:   v.Type(),
		})
	}

	return rf
}
