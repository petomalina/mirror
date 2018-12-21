package bundle

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type PackageSuite struct {
	suite.Suite
}

func (s *PackageSuite) TestListGoFiles() {
	ff, err := listGoFiles(".")
	s.NoError(err)
	s.Len(ff, 5)
}

func TestPackageSuite(t *testing.T) {
	suite.Run(t, &PackageSuite{})
}
