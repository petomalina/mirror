package mirror

import (
	"reflect"
	"strings"
)

// Struct is a wrapper for the runtime symbol Ref
type Struct struct {
	Ref interface{}

	// OriginalPackage can be overridden by package copiers to set the original package path
	// if the package has changed due to movement of the Ref
	// This is mainly used by the CreateDefaultApp when it's copying the package over to the cache
	OriginalPackage string
}

// ReflectStruct creates a mirror Struct that enables users
// to use mirror-enhanced reflections
func ReflectStruct(s interface{}) *Struct {
	return &Struct{
		Ref:             s,
		OriginalPackage: "",
	}
}

// Name returns the name of an underlying type for a given reflection.
// Note that this does not return pointer type asterixes nor
// package names
func (s *Struct) Name() string {
	name := reflect.TypeOf(s.Ref).Elem().String()

	// strip the package prefix, as we don't want it explicitly in the name
	pkgStrip := strings.Split(name, ".")
	if len(pkgStrip) > 1 {
		return pkgStrip[1]
	}

	return pkgStrip[0]
}

// PkgPath returns the import path for the current reflection
func (s *Struct) PkgPath() string {
	if s.OriginalPackage != "" {
		return s.OriginalPackage
	}

	return reflect.TypeOf(s.Ref).Elem().PkgPath()
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

// RawStructFieldType is returned by the RawFields as an array of
// key-value pairs (as map can't be indexed by full structure and
// reflect.StuctField is used without pointers across reflect)
type RawStructFieldType struct {
	Value reflect.Value
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

	sValue := reflect.ValueOf(s.Ref).Elem()

	num := sValue.NumField()
	for i := 0; i < num; i++ {
		v := sValue.Field(i)

		rf = append(rf, RawStructFieldType{
			Value: sValue,
			Field: sValue.Type().Field(i),
			Typ:   v.Type(),
		})
	}

	return rf
}

// StructSlice is a slice of pointers to Struct, used to do multi-struct operations
type StructSlice []*Struct

// ReflectStructs is a plural function for ReflectStruct
func ReflectStructs(ss ...interface{}) StructSlice {
	rss := StructSlice{}

	for _, s := range ss {
		rss = append(rss, ReflectStruct(s))
	}

	return rss
}

// StructSliceEachFunc is a type for Each method above StructSlice
type StructSliceEachFunc func(s *Struct)

// Each loops over the StructSlice and calls the callback for each of its elements
// returning itself at the end
func (s StructSlice) Each(cb StructSliceEachFunc) StructSlice {
	for _, st := range s {
		cb(st)
	}

	return s
}

// PkgPath returns the import path for the current reflection
func (s StructSlice) PkgPaths() []string {
	paths := []string{}

	for _, st := range s {
		path := st.PkgPath()
		if path == "" {
			continue
		}

		// just continue if we already have this one
		for _, p := range paths {
			if path == p {
				continue
			}
		}

		paths = append(paths, path)
	}

	return paths
}
