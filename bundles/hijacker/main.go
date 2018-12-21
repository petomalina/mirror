package main

import (
	"bytes"
	"github.com/petomalina/mirror"
	"github.com/petomalina/mirror/bundle"
	"golang.org/x/tools/go/packages"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var typeTemplate = template.Must(template.New("type").Parse(`
type Hijacked{{ .Name }} {{ .Package}}.{{ .Name }}
`))

type TypeTemplateData struct {
	Package string
	Name    string
}

var hiJackedFieldTemplate = template.Must(template.New("hijackedField").Parse(`
func (h *Hijacked{{ .Name }}) Get{{ .FieldName }}() ({{ .FieldType }}, error) {
	return nil, nil
}

func (h *Hijacked{{ .Name }}) Set{{ .FieldName }}(x {{ .FieldType }}) error {
	return nil
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

		mappings := make(map[string]string)
		for _, f := range r.RawFields() {
			// only hijack unexported fields
			if f.Exported() {
				continue
			}

			// TODO: resolve the real type
			mappings[strings.Title(f.Field.Name)] = "string"
		}

		for fieldName, fieldType := range mappings {
			err := hiJackedFieldTemplate.Execute(&blocks, &FieldTemplateData{
				Name:      r.Name(),
				FieldName: fieldName,
				FieldType: fieldType,
			})

			if err != nil {
				return err
			}
		}
	}

	return ioutil.WriteFile(filepath.Join(outDir, "hijacker.go"), blocks.Bytes(), os.ModePerm)
}
