package domain

type User struct {
	ID       string
	Username string
}

func (u User) GetPID() string {
	return u.ID
}

func (u *User) PutPID(pid string) {
	u.ID = pid
}
