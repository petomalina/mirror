package bundle

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type PluginSuite struct {
	suite.Suite
}

func (s *PluginSuite) TestBuildPlugin() {
	//plug, err := copyPackageToCache("./fixtures/user")
	//s.NoError(err)
	//defer func() {
	//	s.NoError(os.RemoveAll(plug))
	//}()
	//
	//so, err := BuildPlugin(plug)
	//s.NoError(err)
	//s.cleanup = append(s.cleanup, so)
	//
	//syms, err := LoadPluginSymbols(so, []string{"XUser"})
	//s.NoError(err)
	//s.Len(syms, 1)
}

func (s *PluginSuite) TestLoadPluginSymbols() {

}

func (s *PluginSuite) TestWithChangedPackage() {

}

func TestPluginSuite(t *testing.T) {
	suite.Run(t, &PluginSuite{})
}
