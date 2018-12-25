package bundle

import (
	"bytes"
	"golang.org/x/tools/go/packages"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

// Writer is an encapsulation of methods used to write to the output directory
// which must represent the package
type Writer struct {
	pkgPath string
	Files   map[string]*File
}

// NewWriter creates a new Writer wrapper
func NewWriter(pkgPath string) *Writer {
	return &Writer{
		pkgPath: pkgPath,

		Files: make(map[string]*File),
	}
}

// File returns an existing file reference or creates a new one which can be manipulated
func (o *Writer) File(name string) *File {
	if f, ok := o.Files[name]; ok {
		return f
	}

	f := &File{
		path: filepath.Join(o.pkgPath, name),
		buf:  &bytes.Buffer{},
	}

	o.Files[name] = f

	return f
}

type File struct {
	path string
	buf  *bytes.Buffer

	Imports []string
}

func (f *File) AddStringTemplate(str string, data interface{}) error {
	return template.Must(template.New(f.path).Parse(str)).Execute(f.buf, data)
}

func (f *File) AddTemplate(t *template.Template, data interface{}) error {
	return t.Execute(f.buf, data)
}

func (f *File) Write() error {
	pkgName, err := DeterminePackage(filepath.Dir(f.path))
	if err != nil {
		return err
	}

	headerBuf := &bytes.Buffer{}
	headerBuf.Write([]byte("package " + pkgName + "\n\n"))

	if len(f.Imports) != 0 {
		headerBuf.WriteString("import (\n")
		for _, i := range f.Imports {
			headerBuf.WriteString("\t\"" + i + "\"\n")
		}
		headerBuf.WriteString(")\n\n")
	}

	return ioutil.WriteFile(f.path, append(headerBuf.Bytes(), f.buf.Bytes()...), os.ModePerm)
}

// DeterminePackage returns a package name for the given directory
// if no package exists, the directory name will be used instead
func DeterminePackage(pkgPath string) (string, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode:  packages.LoadFiles,
		Tests: false,
	}, pkgPath)

	if err != nil || len(pkgs) == 0 || pkgs[0].Name == "" {
		abs, err := filepath.Abs(pkgPath)
		if err != nil {
			return "", err
		}

		return filepath.Base(abs), nil
	}

	return pkgs[0].Name, nil
}
