package plugins

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

type LoaderSuite struct {
	suite.Suite
}

type LoaderCandidate struct {
	name    string
	loader  *Loader
	symbols []string

	err            error
	len            int
	preservedCache bool
}

type WatchCandidate struct {
	name    string
	loader  *Loader
	symbols []string

	errs             []error
	loadedSymbolsLen int
	triggerChange    bool
}

const (
	TestCacheDir     = ".testmirror"
	ReadonlyCacheDir = ".testmirror-readonly"
)

func (s *LoaderSuite) SetupTest() {
	s.NoError(os.Mkdir(ReadonlyCacheDir, 0400))
}

func (s *LoaderSuite) TearDownTest() {
	s.CleanupCacheDirs()

	// clean the readonly cache after the test as we won't
	// shouldn't be able to copy anything to that folder anyway
	s.NoError(os.Remove(ReadonlyCacheDir))
}

func (s *LoaderSuite) CleanupCacheDirs() {
	// remove the cache dir if it exists
	s.NoError(os.RemoveAll(DefaultCache))
	s.NoError(os.RemoveAll(TestCacheDir))
}

func (s *LoaderSuite) TestLoad() {
	candidates := []LoaderCandidate{
		{
			name: "Load already created symbols from ./fixtures/user",
			loader: &Loader{
				TargetPath: "./fixtures/user",
			},
			symbols: []string{"XUser"},
			len:     1,
		},
		{
			name: "Load already created symbols using 'all' from ./fixtures/user",
			loader: &Loader{
				TargetPath: "./fixtures/user",
			},
			symbols: []string{"all"},
			len:     1,
		},
		{
			name: "Generate symbols for ./fixtures/usernosymbol",
			loader: &Loader{
				TargetPath:      "./fixtures/usernosymbol",
				GenerateSymbols: true,
			},
			symbols: []string{"User"},
			len:     1,
		},
		{
			name: "Preserve cache for ./fixtures/user",
			loader: &Loader{
				TargetPath:    "./fixtures/user",
				PreserveCache: true,
			},
			symbols:        []string{"XUser"},
			len:            1,
			preservedCache: true,
		},
		{
			name: "Preserve cache in a different dir for ./fixtures/user",
			loader: &Loader{
				TargetPath:    "./fixtures/user",
				PreserveCache: true,
				CacheDir:      TestCacheDir,
			},
			symbols:        []string{"XUser"},
			len:            1,
			preservedCache: true,
		},
		{
			name: "Get error when copying to root from ./fixtures/usernosymbol",
			loader: &Loader{
				TargetPath: "./fixtures/usernosymbol",
				CacheDir:   ReadonlyCacheDir,
			},
			err: ErrCopyingToCacheFailed,
		},
		{
			name: "Get non-existing package error for ./fixtures/nonexisting",
			loader: &Loader{
				TargetPath: "./fixtures/nonexisting",
			},
			err: ErrFindPackageFailed,
		},
		{
			name: "Get symbol load fail for XUser without symbol generation in ./fixtures/usernosymbol",
			loader: &Loader{
				TargetPath: "./fixtures/usernosymbol",
			},
			symbols: []string{"XUser"},
			err:     ErrSymbolLoadFailed,
		},
	}

	for _, c := range candidates {
		fmt.Println("Running test case:", c.name)

		model, err := c.loader.Load(c.symbols)

		s.EqualValues(c.err, errors.Cause(err))
		s.Len(model, c.len)

		if c.preservedCache {
			actualCacheDir := DefaultCache
			if c.loader.CacheDir != "" {
				actualCacheDir = c.loader.CacheDir
			}

			ff, err := ioutil.ReadDir(actualCacheDir)
			s.NoError(err, "cachedir %s not found", actualCacheDir)
			s.NotEqual(0, len(ff))
		}

		s.CleanupCacheDirs()
	}
}

func (s *LoaderSuite) TestWatch() {
	candidates := []WatchCandidate{
		{
			name: "Watch only done for created symbols in ./fixtures/user",
			loader: &Loader{
				TargetPath: "./fixtures/user",
			},
			symbols: []string{"XUser"},
		},
		{
			name: "Watch already created symbols in ./fixtures/user (triggered)",
			loader: &Loader{
				TargetPath: "./fixtures/user",
			},
			symbols:          []string{"XUser"},
			loadedSymbolsLen: 1,
			triggerChange:    true,
		},
		{
			name: "Get errors within ./fixtures/usernosymbol without generating (triggered)",
			loader: &Loader{
				TargetPath: "./fixtures/usernosymbol",
			},
			symbols:       []string{"XUser"},
			triggerChange: true,
			errs:          []error{ErrSymbolLoadFailed},
		},
		{
			name: "Get error when copying to root from ./fixtures/usernosymbol",
			loader: &Loader{
				TargetPath: "./fixtures/usernosymbol",
				CacheDir:   ReadonlyCacheDir,
			},
			errs:          []error{ErrCopyingToCacheFailed},
			triggerChange: true,
		},
	}

	for _, c := range candidates {
		fmt.Println("Running test case:", c.name)

		done := make(chan bool)
		modelsChan, errChan := c.loader.Watch(c.symbols, done)

		s.NotNil(modelsChan)
		s.NotNil(errChan)

		// create and remove a file in the folder to triggerChange a random change
		if c.triggerChange {
			fName := filepath.Join(c.loader.TargetPath, fmt.Sprintf("%d.txt", rand.Int()))
			s.NoError(ioutil.WriteFile(fName, []byte("hello world"), os.ModePerm))
			s.NoError(os.Remove(fName))
		}

		// setup helper for expected number of triggers so we can break the channels
		expectedModelTriggersCount := 0
		if c.triggerChange && len(c.errs) == 0 {
			expectedModelTriggersCount = 1
		}

		errTriggerCounter := 0
		modelTriggerCounter := 0
		for {
			// break out of the loop if we received everything
			if modelTriggerCounter == expectedModelTriggersCount && errTriggerCounter == len(c.errs) {
				break
			}

			select {
			case err, ok := <-errChan:
				// check if we found all errors
				if !ok {
					s.EqualValues(len(c.errs), errTriggerCounter)
					break
				}

				s.EqualValues(c.errs[errTriggerCounter], errors.Cause(err))
				errTriggerCounter++
			case syms, ok := <-modelsChan:
				if !ok {
					s.EqualValues(expectedModelTriggersCount, modelTriggerCounter)
					break
				}

				s.EqualValues(c.loadedSymbolsLen, len(syms))
				modelTriggerCounter++
			}
		}

		done <- true
	}
}

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, &LoaderSuite{})
}
