package main

import (
	"github.com/petomalina/mirror"
	"golang.org/x/tools/go/packages"
	"log"
	"path/filepath"
	"strings"
	"text/template"
)

var functionalTemplate = template.Must(template.New("funcional").Parse(
	`type _T_Slice []*_T_
type _T_MapCallback func(*_T_) *_T_

// Map replaces each object in slice by its mapped descendant
func (us _T_Slice) Map(cb _T_MapCallback) _T_Slice {
	newSlice := _T_Slice{}
	for _, o := range us {
		newSlice = append(newSlice, cb(o))
	}
	
	return newSlice
}

type _T_FilterCallback func(*_T_) bool

func (us _T_Slice) Filter(cb _T_FilterCallback) _T_Slice {
	newSlice := _T_Slice{}
	for _, o := range us  {
		if cb(o) {
			newSlice = append(newSlice, o)
		}
	}

	return newSlice
}

type _T_ReduceCallback func(interface{}, *_T_) interface{} 

func (us _T_Slice) Reduce(cb _T_ReduceCallback, init interface{}) interface{} {
	var res = init
	for _, o := range us {
		res = cb(res, o)
	}

	return res
}
`))

func main() {
	if err := mirror.RunDefaultApp("mirror-functional", ProcessModel); err != nil {
		log.Fatal(err)
	}
}

func ProcessModel(outDir string, models []interface{}, _ *packages.Package) error {
	out := mirror.File(filepath.Join(outDir, "functional.go"))
	blocks := []string{}

	for _, rs := range mirror.ReflectStructs(models...) {
		blocks = append(blocks, strings.Replace(functionalTemplate, "_T_", rs.Name(), -1))
	}

	return out.Write(blocks...)
}
