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

	AssertReflectedStruct(
		&s.Suite,
		expectedReflection{
			name: "User",
			pkg:  "github.com/petomalina/mirror/fixtures/user",
		},
		ReflectStruct(userFixture.XUser),
	)
}

type expectedReflection struct {
	name string
	pkg  string
}

func AssertReflectedStruct(s *suite.Suite, ex expectedReflection, ref *Struct) {
	s.EqualValues(ex.name, ref.Name())
	s.EqualValues(ex.pkg, ref.PkgPath())
}

func (s *StructSuite) TestRawFields() {
	fields := ReflectStruct(userFixture.XUser).RawFields()

	s.EqualValues(3, len(fields))
	s.EqualValues("Email", fields[0].Field.Name)
	s.EqualValues("string", fields[0].Typ.String())
	s.EqualValues("Name", fields[1].Field.Name)
	s.EqualValues("string", fields[1].Typ.String())
	s.EqualValues("password", fields[2].Field.Name)
	s.EqualValues("string", fields[2].Typ.String())
}

func (s *StructSuite) TestRawFieldExported() {
	fields := ReflectStruct(userFixture.XUser).RawFields()
	s.True(fields[0].Exported())
	s.True(fields[1].Exported())
	s.False(fields[2].Exported())
}

func TestStructSuite(t *testing.T) {
	suite.Run(t, &StructSuite{})
}
