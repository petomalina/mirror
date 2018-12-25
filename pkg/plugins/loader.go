package plugins

import "os"

// Loader encapsulates full lifecycle of a plugin
type Loader struct {
	// TargetPath is a relative path to the plugin that should be built
	TargetPath string

	// PreserveCache determines if the copied cached files should be preserved
	// at the end of the lifecycle
	PreserveCache bool

	// GenerateSymbols automatically generates symbols for desired types,
	// e.g. type User struct {} will have his var XUser User generated
	// in case it is passed into symbolNames for Load function
	GenerateSymbols bool

	// CacheDir can be used to override the default setting of the cache directory
	// if not set, this will default to the DefaultCache
	CacheDir string
}

func (l *Loader) Load(symbolNames []string) ([]interface{}, error) {
	pkg, err := FindPackage(l.TargetPath)
	if err != nil {
		return nil, err
	}

	// copy everything into the cache so we can manipulate it further and avoid caching
	cacheTargetPath, err := CopyPackageToCache(pkg, DefaultCache)
	if err != nil {
		return nil, err
	}

	// generate symbols and accept new symbol names
	if l.GenerateSymbols {
		symbolNames, err = GenerateSymbolsForModels(symbolNames, cacheTargetPath)
		if err != nil {
			return nil, err
		}
	}

	so, err := Build(cacheTargetPath, cacheTargetPath)
	if err != nil {
		return nil, err
	}

	syms, err := LoadSymbols(so, symbolNames)
	if err != nil {
		return nil, err
	}

	// don't do cleanup if we want to preserve the cache, just return
	if l.PreserveCache {
		return syms, nil
	}

	return syms, os.RemoveAll(cacheTargetPath)
}

func (l *Loader) Watch(symbolnames []string) (chan<- []interface{}, error) {
	return nil, nil
}
