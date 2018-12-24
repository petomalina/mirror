package main

import (
	"github.com/petomalina/mirror"
	"golang.org/x/tools/go/packages"
	"log"
	"strings"
	"text/template"
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
	return "", nil
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
	if err := mirror.RunDefaultApp("mirror-hijacker", ProcessModel); err != nil {
		log.Fatal(err)
	}
}

func ProcessModel(models mirror.StructSlice, out *mirror.Out, pkg *packages.Package) error {
	temp := out.File("hijacker.go")
	temp.Imports = models.PkgPaths()

	// hijack and create types with fields
	for _, r := range models {
		err := temp.AddTemplate(typeTemplate, &TypeTemplateData{
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
			err := temp.AddTemplate(hiJackedFieldTemplate, &FieldTemplateData{
				Name:      r.Name(),
				FieldName: fieldName,
				FieldType: fieldType,
			})

			if err != nil {
				return err
			}
		}
	}

	return temp.Write()
}
