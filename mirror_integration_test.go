package mirror

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type StructIntegrationSuite struct {
	suite.Suite

	cleanup []string
}

func (s *StructIntegrationSuite) TearDownTest() {
	if len(s.cleanup) > 0 {
		for _, f := range s.cleanup {
			err := os.Remove(f)
			if err != nil {
				L.Warn("An error occured during cleanup of: ", f, " : ", err.Error())
			}
		}

		s.cleanup = []string{}
	}
}

func (s *StructIntegrationSuite) TestReflectStruct() {
	so, err := BuildPlugin("./fixtures/user")
	s.NoError(err)
	s.cleanup = append(s.cleanup, so)

	//_, err = LoadPluginSymbols(so, []string{"XUser"})
	//s.NoError(err)
	//s.Len(syms, 1)

	//assertReflectedStruct(
	//	&s.Suite,
	//	expectedReflection{
	//		name: "User",
	//	},
	//	ReflectStruct(syms[0]),
	//)
}

func TestStructIntegrationSuite(t *testing.T) {
	//L.SetLevel(logrus.TraceLevel)
	suite.Run(t, &StructIntegrationSuite{})
}
