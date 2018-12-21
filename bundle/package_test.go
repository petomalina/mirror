package bundle

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type PackageSuite struct {
	suite.Suite
}

func (s *PackageSuite) TestListGoFiles() {
	ff, err := findPackage(".")
	s.NoError(err)
	s.Len(ff.GoFiles, 5)
}

func TestPackageSuite(t *testing.T) {
	suite.Run(t, &PackageSuite{})
}
