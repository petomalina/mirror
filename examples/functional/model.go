package functional

//go:generate functional -m XUser

type Email string

type User struct {
	Name  string
	Email Email
}
