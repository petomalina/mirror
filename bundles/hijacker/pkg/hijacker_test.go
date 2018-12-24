package pkg

import (
	"github.com/petomalina/mirror"
	"github.com/petomalina/mirror/fixtures/user"
	"github.com/stretchr/testify/suite"
	"testing"
)

type HijackerSuite struct {
	suite.Suite
}

func (s *HijackerSuite) TestHijack() {
	err := Hijack(mirror.ReflectStruct(user.XUser))
	s.NoError(err)
}

func TestHijackerSuite(t *testing.T) {
	suite.Run(t, &HijackerSuite{})
}
