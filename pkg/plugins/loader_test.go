package plugins

import (
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

func (s *LoaderSuite) TearDownTest() {
	// remove the cache dir if it exists
	s.NoError(os.RemoveAll(DefaultCache))
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
	}

	for _, c := range candidates {
		model, err := c.loader.Load(c.symbols)

		s.EqualValues(c.err, err)
		s.Len(model, c.len)

		if c.preservedCache {
			ff, err := ioutil.ReadDir(DefaultCache)
			s.NoError(err)
			s.NotEqual(0, len(ff))
		}
	}
}

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, &LoaderSuite{})
}
