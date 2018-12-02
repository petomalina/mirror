package main

import (
	"github.com/petomalina/mirror"
	"github.com/petomalina/mirror/examples/functional"
	"strings"
)

const mapTemplate = `type _T_Slice []*_T_

type _T_MapCallback func(*_T_) *_T_

func (us _T_Slice) Map(cb _T_MapCallback) _T_Slice {
	newSlice := _T_Slice{}
	for _, o := range us {
		newSlice = append(newSlice, cb(o))
	}
	
	return newSlice
}
`

func main() {
	user := mirror.ReflectStruct(functional.User{})

	out := mirror.File("functional.go")

	blocks := []string{}
	blocks = append(blocks, strings.Replace(mapTemplate, "_T_", user.Name(), -1))

	//for fName, fType := range user.Fields() {
	//	// if the type is from the same package, trim the package name
	//	// e.g. functional.Email will become Email
	//	fType = strings.TrimPrefix(fType, mirror.Package()+".")
	//
	//	blocks = append(blocks, strings.Replace(mapTemplate, "_T_", fName, -1))
	//}

	err := out.Write(blocks...)
	if err != nil {
		panic(err)
	}
}
