package bundle

import (
	. "github.com/petomalina/mirror/pkg/logger"
	"github.com/petomalina/mirror/pkg/plugins"
	"golang.org/x/tools/go/packages"
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

	loader := plugins.Loader{
		TargetPath:      pkgName,
		GenerateSymbols: opts.GenerateSymbols,
		PreserveCache:   opts.PreserveCache,
	}

	models, err := loader.Load(symbols)

	L.Method("Bundle", "Run").Trace("Loaded symbols: ", models)
	return models, pkg, err

}
