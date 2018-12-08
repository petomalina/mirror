package mirror

import (
	"errors"
	"io/ioutil"
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

// WithChangedPackage changes the `package X` line of each file in the
// targeted package, changing its name to the desiredPkgName, running the
// `run` function and changing it back to the default
func WithChangedPackage(pkg, desiredPkgName string, run func() error) error {
	goFiles, err := listGoFiles(pkg)
	if err != nil {
		return err
	}

	// fileContents map holds all changed files with their
	// original package names and contents
	fileContents, err := readFilesAndPackages(goFiles)
	if err != nil {
		return err
	}

	// replace all package directives to the desired package names
	for _, f := range goFiles {
		err = ioutil.WriteFile(
			f,
			[]byte(pkgRegex.ReplaceAll(fileContents[f].content, []byte("package "+desiredPkgName))),
			0,
		)

		if err != nil {
			return err
		}
	}

	if err = run(); err != nil {
		return err
	}

	// replace back to the original package names
	for _, f := range goFiles {
		err = ioutil.WriteFile(
			f,
			[]byte(pkgRegex.ReplaceAll(
				fileContents[f].content,
				append([]byte("package "), fileContents[f].originalPkgName...),
			)),
			0,
		)

		if err != nil {
			return err
		}
	}

	return nil
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

		ff = append(ff, fi.Name())
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
