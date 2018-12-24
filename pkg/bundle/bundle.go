package bundle

import (
	"fmt"
	"github.com/petomalina/mirror/pkg/cp"
	"github.com/petomalina/mirror/pkg/plugins"
	"golang.org/x/tools/go/packages"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"

	. "github.com/petomalina/mirror/pkg/logger"
)

// Bundle is a set of templates and logic packed for a purpose of
// code generation that is generalized (into the bundle)
type Bundle struct{}

// RunOptions encapsulates options that set settings for the Run method
type RunOptions struct {
	GenerateSymbols bool
	PreserveCache   bool
}

func (b *Bundle) Run(pkgName string, symbols []string, opts RunOptions) ([]interface{}, *packages.Package, error) {
	L.Method("Bundle", "Run").Trace("Invoked, should load: ", symbols, "from: ", pkgName)

	pkg, err := plugins.FindPackage(pkgName)
	if err != nil {
		return nil, nil, err
	}

	// cp to the cache dir
	pkgCacheDir, err := copyPackageToCache(pkg)
	if err != nil {
		return nil, nil, err
	}

	// automatically generate all symbols, creating e.g. XUser from User etc.
	if opts.GenerateSymbols {
		err = generateSymbolsForModels(symbols, pkgCacheDir)
		if err != nil {
			return nil, nil, err
		}

		// mutate to match the symbol prefix
		for i := range symbols {
			symbols[i] = "X" + symbols[i]
		}
	}

	// remove the cache dir once we are done
	// only if user doesn't want to preserve it
	if !opts.PreserveCache {
		L.Method("Bundle", "Run").Trace("Cache will be deleted, as preserveCache is not set")
		defer func(dir string) {
			L.Method("Bundle", "Run").Trace("Removing cache dir: ", dir)
			if err := os.RemoveAll(dir); err != nil {
				L.Method("Bundle", "Run").Warn("An error occurred when removing cache dir: " + dir)
			}
		}(pkgCacheDir)
	}

	L.Method("Bundle", "Run").Trace("Building the plugin: ", pkgCacheDir)
	objPath, err := plugins.Build(pkgCacheDir)
	if err != nil {
		return nil, nil, err
	}

	// remove the object model when we are done
	// this object resides in the .mirror folder, that's why it's not
	// cleared by the cp cleaner
	defer func(oPath string) {
		L.Method("Bundle", "Run").Trace("Removing object model: ", oPath)
		if err := os.Remove(oPath); err != nil {
			L.Method("Bundle", "Run").Warn("An error occurred when removing: " + oPath)
		}
	}(objPath)

	L.Method("Bundle", "Run").Trace("Opening the plugin: ", objPath)
	models, err := plugins.LoadSymbols(objPath, symbols)

	L.Method("Bundle", "Run").Trace("Loaded symbols: ", models)
	return models, pkg, err

}

func generateSymbolsForModels(models []string, out string) error {
	symbolsFile := filepath.Join(out, fmt.Sprintf("/%d.go", rand.Int()))

	tmpl := `// DO NOT EDIT: THIS BLOCK IS AUTOGENERATED BY MIRROR BUNDLE
package main

var (
`

	for _, m := range models {
		tmpl += "\tX" + m + "  = " + m + "{}\n"
	}

	tmpl += ")\n"

	return ioutil.WriteFile(symbolsFile, []byte(tmpl), os.ModePerm)
}

func copyPackageToCache(pkg *packages.Package) (string, error) {
	// Copy the directory so the plugin can be build outside
	pkgCacheDir := fmt.Sprintf("./.mirror/%d", rand.Int())
	L.Method("Bundle", "Run").Trace("Making cache dir: ", pkgCacheDir)
	err := os.MkdirAll(pkgCacheDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	L.Method("Bundle", "Run").Trace("Copying ", pkg, "->", pkgCacheDir)
	for _, f := range pkg.GoFiles {
		err := cp.File(f, filepath.Join(pkgCacheDir, filepath.Base(f)))
		if err != nil {
			return pkgCacheDir, err
		}
	}

	return pkgCacheDir, nil
}
