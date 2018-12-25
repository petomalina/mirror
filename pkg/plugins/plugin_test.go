package plugins

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type PluginSuite struct {
	suite.Suite
}

func (s *PluginSuite) TestBuildPlugin() {}

func (s *PluginSuite) TestLoadPluginSymbols() {

}

func (s *PluginSuite) TestWithChangedPackage() {

}

func TestPluginSuite(t *testing.T) {
	suite.Run(t, &PluginSuite{})
}
