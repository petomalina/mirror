package main

import (
	"bytes"
	"fmt"
	"github.com/petomalina/mirror"
	"github.com/petomalina/mirror/bundle"
	"golang.org/x/tools/go/packages"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var typeTemplate = template.Must(template.New("type").Parse(`
type Hijacked{{ .Name }} *{{ .Package}}.{{ .Name }}
`))

type TypeTemplateData struct {
	Package string
	Name    string
}

var hiJackedFieldTemplate = template.Must(template.New("hijackedField").Parse(`
func (h *Hijacked{{ .Name }}) Set{{ .FieldName }}(x {{ .FieldType }}) error {
	return nil
}

func (h *Hijacked{{ .Name }}) Get{{ .FieldName }}() ({{ .FieldType }}, error) {
	return nil, nil
}
`))

type FieldTemplateData struct {
	Name      string
	FieldName string
	FieldType string
}

func main() {
	b := &bundle.Bundle{
		RunFunc: ProcessModel,
	}

	if err := b.RunDefaultApp("mirror-hijacker"); err != nil {
		log.Fatal(err)
	}
}

func ProcessModel(outDir string, models []interface{}, pkg *packages.Package) error {
	blocks := bytes.Buffer{}

	absPath, _ := filepath.Abs(outDir)
	// write the package and the import for the original package
	blocks.Write([]byte("package " + filepath.Base(absPath) + "\n\nimport \"" + pkg.PkgPath + "\"\n"))

	// reflect all models
	rs := mirror.ReflectStructs(models...)

	// hijack and create types with fields
	for _, r := range rs {
		err := typeTemplate.Execute(&blocks, &TypeTemplateData{
			Package: pkg.Name,
			Name:    r.Name(),
		})
		if err != nil {
			return err
		}

		for _, f := range r.RawFields() {
			fmt.Println(f)
		}
	}

	return ioutil.WriteFile(filepath.Join(outDir, "hijacker.go"), blocks.Bytes(), os.ModePerm)
}
