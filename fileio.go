package mirror

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// File represents a path to the output file
type File string

func (f File) Package() string {
	wd, _ := os.Getwd()

	pkg := PackageFromPath(filepath.Join(wd, string(f)))

	return fmt.Sprintf("package %s\n\n", pkg)
}

// Write creates the necessary path in the filesystem, converts
// given strings into bytes and writes them into the output file
func (f File) Write(content ...string) error {
	bb := bytes.Buffer{}
	bb.WriteString(f.Package())

	for _, c := range content {
		// Sprintln also leaves NewLine at the end of the string
		_, err := bb.WriteString(fmt.Sprintln(c))
		if err != nil {
			return err
		}
	}

	err := os.MkdirAll(filepath.Dir(string(f)), os.ModePerm)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(string(f), bb.Bytes(), os.ModePerm)
}
