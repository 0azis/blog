package domain

type User struct {
	ID          int     `json:"-" db:"id"`
	Email       string  `json:"email" db:"email"`
	Username    string  `json:"username" db:"username"`
	Password    string  `json:"-" db:"password"`
	Name        *string `json:"name" db:"name"`
	Avatar      *string `json:"avatar" db:"avatar"`
	Description *string `json:"description" db:"description"`
	Owner       bool    `json:"owner"`
}

func (u *User) SetOwnership(jwtUserID int) {
	if u.ID == jwtUserID {
		u.Owner = true
	}
}

type SignUpCredentials struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
