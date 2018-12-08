package functional

//go:generate MIRROR_MODELS=USer go run ./mirror/generate.go

type Email string

type User struct {
	Name  string
	Email Email
}

// GENERATOR ONLY
// These variables are only exported for the purpose
// of code generation
var (
	XUser = User{}
)
