package domain

import (
	"regexp"
)

type User struct {
	ID          int     `json:"id" db:"id"`
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

type Author struct {
	ID       int     `json:"id" db:"id"`
	Username string  `json:"username" db:"username"`
	Name     string  `json:"name" db:"name"`
	Avatar   *string `json:"avatar" db:"avatar"`
}

type UserPatch struct {
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
}

func (up UserPatch) Validate() bool {
	return up.Name != "" && up.Avatar != "" && up.Description != ""
}

type ValidationalData struct {
	IsEmail    bool `json:"isEmail"`
	IsUsername bool `json:"isUsername"`
	IsPassword bool `json:"isPassword"`
}

func ValidateUser(credentials SignUpCredentials, user User) ValidationalData {
	return ValidationalData{
		IsEmail:    user.Email != credentials.Email,
		IsUsername: user.Username != credentials.Username,
		IsPassword: goodPassword(credentials.Password),
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

func goodPassword(password string) bool {
	containBase, _ := regexp.Match(`[a-z0-9]`, []byte(password))
	containUpper, _ := regexp.Match(`[A-Z]`, []byte(password))
	containSymbols, _ := regexp.Match(`[!@#$%^&*_-]`, []byte(password))

	return (len(password) > 8 && len(password) < 20) && containBase && containUpper && containSymbols
}
