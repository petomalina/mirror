package bundle

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"

	. "github.com/petomalina/mirror/pkg/logger"
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

// TODO: this integration test should be moved into the bundle package
func (s *StructIntegrationSuite) TestReflectStruct() {
	b := Bundle{}

	models, _, err := b.Run("./fixtures/user", []string{"XUser"}, RunOptions{})
	s.NoError(err)
	s.Len(models, 1)
}

func TestStructIntegrationSuite(t *testing.T) {
	//L.SetLevel(logrus.TraceLevel)
	suite.Run(t, &StructIntegrationSuite{})
}
