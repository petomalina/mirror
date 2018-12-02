package main

import (
	"github.com/petomalina/mirror"
	"log"
	"path/filepath"
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
	err := mirror.RunDefaultApp("mirror-functional", &mirror.Bundle{
		RunFunc: ProcessModel,
	})

	if err != nil {
		log.Fatal(err)
	}
}

func ProcessModel(outDir string, models []interface{}) error {
	out := mirror.File(filepath.Join(outDir, "functional.go"))
	blocks := []string{}

	for _, m := range models {
		str := mirror.ReflectStruct(m)

		blocks = append(blocks, strings.Replace(mapTemplate, "_T_", str.Name(), -1))
	}

	return out.Write(blocks...)
}
