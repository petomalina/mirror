package mirror

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"plugin"
)

// BuildPlugin builds the given package into plugin and saves it in
// current path under a random name .so, returning the name to the caller
func BuildPlugin(pkg string) (string, error) {
	// random file name so we'll get unique loader each time
	uniq := rand.Int()

	objPath := fmt.Sprintf("mirror-%d.so", uniq)
	L.Method("Bundle", "Run").Trace("Object path: ", objPath)

	// create the plugin from the passed package
	err := WithChangedPackage(pkg, "main", func() error {
		cmd := exec.Command("go", "build", "-buildmode=plugin", "-o="+objPath, pkg)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	})

	return objPath, err
}

// LoadPluginSymbols accepts a plugin path and returns all symbols
// that were found in the given plugin.
// If * is provided as only value in `symbols`, all symbols from the
// given plugin will be returned
func LoadPluginSymbols(pluginPath string, symbols []string) ([]interface{}, error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}

	// special case - load all exported symbols from the file
	if len(symbols) == 1 && symbols[0] == "*" {
		// TODO: make it possible via reflection (get all symbols from plugin)
	}

	// add model symbols that were loaded from the built plugin
	models := []interface{}{}
	for _, symName := range symbols {
		sym, err := p.Lookup(symName)
		if err != nil {
			return nil, err
		}

		models = append(models, sym)
	}

	return models, nil
}

// WithChangedPackage changes the `package X` line of each file in the
// targeted package, changing its name to the desiredPkgName, running the
// `run` function and changing it back to the default
func WithChangedPackage(pkg, desiredPkgName string, run func() error) error {
	L.Method("Internal/package", "WithChangedPackage").Trace("Invoked on pkg: ", pkg)
	goFiles, err := listGoFiles(pkg)
	if err != nil {
		return err
	}
	L.Method("Internal/package", "WithChangedPackage").Trace("Listed files: ", goFiles)

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
