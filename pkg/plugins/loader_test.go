package plugins

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"testing"
)

type LoaderSuite struct {
	suite.Suite
}

type LoaderCandidate struct {
	loader  *Loader
	symbols []string

	err            error
	len            int
	preservedCache bool
}

const (
	TestCacheDir = ".testmirror"
)

func (s *LoaderSuite) TearDownTest() {
	s.CleanupCacheDirs()
}

func (s *LoaderSuite) CleanupCacheDirs() {
	// remove the cache dir if it exists
	s.NoError(os.RemoveAll(DefaultCache))
	s.NoError(os.RemoveAll(TestCacheDir))
}

func (s *LoaderSuite) TestLoad() {
	candidates := []LoaderCandidate{
		{
			loader: &Loader{
				TargetPath: "./fixtures/user",
			},
			symbols: []string{"XUser"},
			len:     1,
		},
		{
			loader: &Loader{
				TargetPath:      "./fixtures/usernosymbol",
				GenerateSymbols: true,
			},
			symbols: []string{"User"},
			len:     1,
		},
		{
			loader: &Loader{
				TargetPath:    "./fixtures/user",
				PreserveCache: true,
			},
			symbols:        []string{"XUser"},
			len:            1,
			preservedCache: true,
		},
		{
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
			loader: &Loader{
				TargetPath: "./fixtures/usernosymbol",
				CacheDir:   "/.mirror",
			},
			err: ErrCopyingToCacheFailed,
		},
		{
			loader: &Loader{
				TargetPath: "./fixtures/nonexisting",
			},
			err: ErrFindPackageFailed,
		},
		{
			loader: &Loader{
				TargetPath: "./fixtures/usernosymbol",
			},
			symbols: []string{"XUser"},
			err:     ErrSymbolLoadFailed,
		},
	}

	for _, c := range candidates {
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

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, &LoaderSuite{})
}
