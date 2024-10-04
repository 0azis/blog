package domain

type User struct {
	ID          int     `json:"-" db:"id"`
	FirstName   string  `json:"firstName" db:"first_name"`
	LastName    string  `json:"lastName" db:"last_name"`
	Username    string  `json:"username" db:"username"`
	Password    string  `json:"-" db:"password"`
	Avatar      *string `json:"avatar" db:"avatar"`
	Description *string `json:"description" db:"description"`
}

func (u User) FullName() string {
	return u.FirstName + " " + u.LastName
}

type SignUpCredentials struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type SignInCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
