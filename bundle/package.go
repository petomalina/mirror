package bundle

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

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

// readFilesAndPackages accepts a set of filtered Go files,
// reads them and caches their contents with their package names
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

func copyPackageToCache(pkg string) (string, error) {
	// Copy the directory so the plugin can be build outside
	pkgCacheDir := fmt.Sprintf("./.mirror/%d", rand.Int())
	L.Method("Bundle", "Run").Trace("Making cache dir: ", pkgCacheDir)
	err := os.MkdirAll(pkgCacheDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	L.Method("Bundle", "Run").Trace("Copying ", pkg, "->", pkgCacheDir)
	return pkgCacheDir, CopyDir(pkg, pkgCacheDir, false)
}

func generateSymbolsForModels(models []string, out string) error {
	symbolsFile := filepath.Join(out, fmt.Sprintf("/%d.go", rand.Int()))

	tmpl := `// DO NOT EDIT: THIS BLOCK IS AUTOGENERATED BY MIRROR BUNDLE
package main

var (
`

	for _, m := range models {
		tmpl += "\tX" + m + "  = &" + m + "{}\n"
	}

	tmpl += ")\n"

	return ioutil.WriteFile(symbolsFile, []byte(tmpl), os.ModePerm)
}
