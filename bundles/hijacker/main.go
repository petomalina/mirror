package main

import (
	"github.com/petomalina/mirror"
	"github.com/petomalina/mirror/bundle"
	"log"
	"path/filepath"
	"strings"
)

const mapTemplate = ``

func main() {
	b := &bundle.Bundle{
		RunFunc: ProcessModel,
	}

	if err := b.RunDefaultApp("mirror-hijacker"); err != nil {
		log.Fatal(err)
	}
}

func ProcessModel(outDir string, models []interface{}) error {
	out := mirror.File(filepath.Join(outDir, "functional.go"))
	blocks := []string{}

	for _, rs := range mirror.ReflectStructs(models...) {
		blocks = append(blocks, strings.Replace(mapTemplate, "_T_", rs.Name(), -1))
	}

	return out.Write(blocks...)
}
