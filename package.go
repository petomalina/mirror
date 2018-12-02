package mirror

import (
	"os"
	"path/filepath"
)

// Package returns current package that the generator was run in
func Package() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return PackageFromPath(dir)
}

// PackageFromPath returns the directory name representing package name
// for the given path
func PackageFromPath(dir string) string {
	return filepath.Base(filepath.Dir(dir))
}
