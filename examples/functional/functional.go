package functional

type UserSlice []*User

type UserMapCallback func(*User) *User

// Map replaces each object in slice by its mapped descendant
func (us UserSlice) Map(cb UserMapCallback) UserSlice {
	newSlice := UserSlice{}
	for _, o := range us {
		newSlice = append(newSlice, cb(o))
	}

	return newSlice
}
