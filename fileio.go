package mirror

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
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

// REFERENCE: https://blog.depado.eu/post/copy-files-and-directories-in-go
// CopyFile copies a single file from src to dst
func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// CopyDir copies a whole directory recursively
func CopyDir(src string, dst string, r bool) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			// skip if not recursive
			if !r {
				continue
			}

			if err = CopyDir(srcfp, dstfp, r); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
