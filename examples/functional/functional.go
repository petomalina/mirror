package functional

type UserSlice []*User

type UserMapCallback func(*User) *User

func (us UserSlice) Map(cb UserMapCallback) UserSlice {
	newSlice := UserSlice{}
	for _, o := range us {
		newSlice = append(newSlice, cb(o))
	}
	
	return newSlice
}