package mirror

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

var (
	pkgRegex = regexp.MustCompile(`(?m:^package (?P<pkg>\w+$))`)
)

// changedFileContent is a utility structure to save the package
// change and the contents so we won't need to reload it
type changedFileContent struct {
	content         []byte
	originalPkgName []byte
}

// listGoFiles returns names of go files in the targeted directory
func listGoFiles(dir string) ([]string, error) {
	fii, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	ff := []string{}
	for _, fi := range fii {
		// skip directories as they may contain other packages
		if fi.IsDir() {
			continue
		}

		if !strings.HasSuffix(fi.Name(), ".go") {
			continue
		}

		ff = append(ff, filepath.Join(dir, fi.Name()))
	}

	return ff, nil
}

func readFilesAndPackages(ff []string) (map[string]changedFileContent, error) {
	changes := map[string]changedFileContent{}

	// go through all files, read them and save their original package names
	for _, f := range ff {
		content, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}

		match := pkgRegex.FindStringSubmatch(string(content))
		if len(match) <= 0 {
			// TODO: wrap the error so it can be easily tracked by callers
			return nil, errors.New("No `package` was found when scanning: " + f)
		}

		originalName := match[0]
		// TODO: what is happening here?
		if len(match) > 1 {
			originalName = match[1]
		}

		// save the change so we can revert
		changes[f] = changedFileContent{
			content:         content,
			originalPkgName: []byte(originalName),
		}
	}

	return changes, nil
}
