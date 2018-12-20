package functional

//go:generate functional -m XUser

type Email string

type User struct {
	Name  string
	Email Email
}

// GENERATOR ONLY
// These variables are only exported for the purpose
// of code generation
var (
	XUser = &User{}
)
