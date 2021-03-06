package plugins

import (
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"os"
	"time"
)

var (
	ErrFindPackageFailed      = errors.New("An error occurred when loading plugin")
	ErrCopyingToCacheFailed   = errors.New("Failed to copy package to the cache")
	ErrSymbolGenerationFailed = errors.New("Failed to generate symbols")
	ErrBuildFailed            = errors.New("Failed to build the plugin")
	ErrSymbolLoadFailed       = errors.New("Failed to load symbols from the plugin")
)

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
		return nil, errors.Wrap(ErrFindPackageFailed, err.Error())
	}

	// copy everything into the cache so we can manipulate it further and avoid caching
	cacheDir := DefaultCache
	if l.CacheDir != "" {
		cacheDir = l.CacheDir
	}

	cacheTargetPath, err := CopyPackageToCache(pkg, cacheDir)
	if err != nil {
		return nil, errors.Wrap(ErrCopyingToCacheFailed, err.Error())
	}

	// generate symbols and accept new symbol names
	if l.GenerateSymbols {
		symbolNames, err = GenerateSymbolsForModels(symbolNames, cacheTargetPath)
		if err != nil {
			return nil, errors.Wrap(ErrSymbolGenerationFailed, err.Error())
		}
	}

	so, err := Build(cacheTargetPath, cacheTargetPath)
	if err != nil {
		return nil, errors.Wrap(ErrBuildFailed, err.Error())
	}

	syms, err := LoadSymbols(so, symbolNames)
	if err != nil {
		return nil, errors.Wrap(ErrSymbolLoadFailed, err.Error())
	}

	// don't do cleanup if we want to preserve the cache, just return
	if l.PreserveCache {
		return syms, nil
	}

	return syms, os.RemoveAll(cacheTargetPath)
}

// Watch watches for changes in the given plugin and emits new symbols
// on changes
func (l *Loader) Watch(symbolNames []string, done <-chan bool) (<-chan []interface{}, <-chan error) {
	out := make(chan []interface{})
	errOut := make(chan error)

	// initialize the watcher with the plugin path
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		errOut <- err
		return out, errOut
	}
	err = watcher.Add(l.TargetPath)
	if err != nil {
		errOut <- err
		return out, errOut
	}

	go func(watcher *fsnotify.Watcher, symbolNames []string, out chan []interface{}, done <-chan bool, errs chan error) {
		defer func() {
			close(out)
			close(errs)
		}()

		// create a debounce channel so we will only triggerChange once for multiple changes
		// this is closed automatically by the debouncer when new value is returned, so it must be
		// renewed when new value should be accepted
		events := eventDebounce(time.Millisecond*300, watcher.Events)

	eventLoop:
		for {
			select {
			// look for events on the filesystem
			case _, ok := <-events:
				// this will break only when we renew the debouncer and the watcher is closed
				if !ok {
					break eventLoop
				}

				// renew the debouncer for the new value
				events = eventDebounce(time.Millisecond*300, watcher.Events)

				syms, err := l.Load(symbolNames)
				if err != nil {
					errs <- err
					break
				}

				// distribute the symbols loaded from the plugin
				out <- syms

				// errors can be handled gracefully and should be proxied to the
				// caller, they can then use `done` channel to stop the execution
			case err, ok := <-watcher.Errors:
				if !ok {
					break eventLoop
				}
				if err != nil {
					errs <- err
				}

			case <-done:
				err = watcher.Close()
				if err != nil {
					errs <- err
				}
				break eventLoop

			}
		}
	}(watcher, symbolNames, out, done, errOut)

	return out, errOut
}

// eventDebounce creates a new debouncer that will be closed once the value is returned
func eventDebounce(interval time.Duration, input <-chan fsnotify.Event) <-chan fsnotify.Event {
	out := make(chan fsnotify.Event)

	go func(interval time.Duration, in <-chan fsnotify.Event, out chan fsnotify.Event) {
		var ev fsnotify.Event
		defer func() {
			close(out)
		}()

		// first iteration out of the after interval, so we will not loop in the time.After without value
		item, ok := <-input
		if !ok {
			return
		}
		ev = item

		for {
			select {
			case item, ok := <-input:
				// if we close the input, we'll just close the output as well
				if !ok {
					return
				}
				ev = item
			case <-time.After(interval):
				out <- ev
				return
			}
		}
	}(interval, input, out)

	return out
}
