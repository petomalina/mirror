package user

type User struct {
	Email string
	Name  string

	// unexported password
	password string
}

// GENERATOR ONLY, DON'T USE
var (
	XUser = &User{}
)
