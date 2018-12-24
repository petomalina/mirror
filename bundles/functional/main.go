package main

import (
	"github.com/petomalina/mirror"
	"golang.org/x/tools/go/packages"
	"log"
)

const functionalTemplate = `type {{ .TypeName }}Slice []*{{ .TypeName }}
type {{ .TypeName }}MapCallback func(*{{ .TypeName }}) *{{ .TypeName }}

// Map replaces each object in slice by its mapped descendant
func (us {{ .TypeName }}Slice) Map(cb {{ .TypeName }}MapCallback) {{ .TypeName }}Slice {
	newSlice := {{ .TypeName }}Slice{}
	for _, o := range us {
		newSlice = append(newSlice, cb(o))
	}
	
	return newSlice
}

type {{ .TypeName }}FilterCallback func(*{{ .TypeName }}) bool

func (us {{ .TypeName }}Slice) Filter(cb {{ .TypeName }}FilterCallback) {{ .TypeName }}Slice {
	newSlice := {{ .TypeName }}Slice{}
	for _, o := range us  {
		if cb(o) {
			newSlice = append(newSlice, o)
		}
	}

	return newSlice
}

type {{ .TypeName }}ReduceCallback func(interface{}, *{{ .TypeName }}) interface{} 

func (us {{ .TypeName }}Slice) Reduce(cb {{ .TypeName }}ReduceCallback, init interface{}) interface{} {
	var res = init
	for _, o := range us {
		res = cb(res, o)
	}

	return res
}
`

func main() {
	if err := mirror.RunDefaultApp("mirror-functional", ProcessModel); err != nil {
		log.Fatal(err)
	}
}

type templateData struct {
	TypeName string
}

func ProcessModel(models mirror.StructSlice, out *mirror.Out, _ *packages.Package) error {
	temp := out.File("functional.go")

	for _, rs := range models {
		err := temp.AddStringTemplate(functionalTemplate, &templateData{
			TypeName: rs.Name(),
		})
		if err != nil {
			return err
		}
	}

	return temp.Write()
}
