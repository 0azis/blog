package domain

type User struct {
	ID          int
	FirstName   string
	LastName    string
	Nick        string
	Password    string
	Avatar      string
	Description string
}

func (u User) FullName() string {
	return u.FirstName + " " + u.LastName
}

type SignUpCredentials struct {
	LastName  string
	FirstName string
	Nick      string
	Password  string
}

type SignInCredentials struct {
	Nick     string
	Password string
}
