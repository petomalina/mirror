package bundle

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"plugin"
	"reflect"
	"unsafe"
)

// BuildPlugin builds the given package into plugin and saves it in
// current path under a random name .so, returning the name to the caller
func BuildPlugin(pkg string) (string, error) {
	L.Method("Internal/plugin", "BuildPlugin").Trace("Invoked with pkg: ", pkg)
	// random file name so we'll get unique loader each time
	uniq := rand.Int()

	objPath := fmt.Sprintf(".mirror/mirror-%d.so", uniq)
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
	if len(symbols) == 1 && symbols[0] == "all" {
		L.Method("Internal/plugin", "LoadPluginSymbols").Trace("Got 'all' option, finding symbols")
		// clear the symbols array so it doesn't contain the *
		symbols = []string{}

		// create the reflection from the 'syms' field of Plugin
		symsField := reflect.ValueOf(p).Elem().FieldByName("syms")
		// create an unsafe pointer so we can access that field (disables the runtime protection)
		symsFieldPtr := reflect.NewAt(symsField.Type(), unsafe.Pointer(symsField.UnsafeAddr())).Elem()

		// range through the map and create the symbols in our array
		for sym := range symsFieldPtr.Interface().(map[string]interface{}) {
			symbols = append(symbols, sym)
		}
	}

	L.Method("Internal/plugin", "LoadPluginSymbols").Trace("Looking up symbols: ", symbols)
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
func WithChangedPackage(pkgName, desiredPkgName string, run func() error) error {
	L.Method("Internal/package", "WithChangedPackage").Trace("Invoked on pkgName: ", pkgName)
	pkg, err := listGoFiles(pkgName)
	if err != nil {
		return err
	}
	L.Method("Internal/package", "WithChangedPackage").Trace("Listed package: ", pkg)

	// replace all package directives to the desired package names
	for _, f := range pkg.GoFiles {
		// read the go file first
		bb, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}

		// write it back with the replaced package (into the out dir)
		err = ioutil.WriteFile(
			f,
			[]byte(pkgRegex.ReplaceAll(bb, []byte("package "+desiredPkgName))),
			0,
		)

		if err != nil {
			return err
		}
	}

	L.Method("Internal/package", "WithChangedPackage").Trace("Running the enclosed function")
	return run()
}
