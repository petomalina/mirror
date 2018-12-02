package functional

//go:generate go run ./mirror/generate.go

type Email string

type User struct {
	Name  string
	Email Email
}
