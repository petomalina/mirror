package mirror

import (
	userFixture "github.com/petomalina/mirror/fixtures/user"
	"github.com/stretchr/testify/suite"
	"testing"
)

type StructSuite struct {
	suite.Suite
}

func (s *StructSuite) TestReflectStruct() {

	assertReflectedStruct(
		&s.Suite,
		expectedReflection{
			name: "User",
		},
		ReflectStruct(userFixture.User{}),
	)
}

type expectedReflection struct {
	name string
}

func assertReflectedStruct(s *suite.Suite, ex expectedReflection, ref *Struct) {
	s.EqualValues(ex.name, ref.Name())
}

func TestStructSuite(t *testing.T) {
	suite.Run(t, &StructSuite{})
}
