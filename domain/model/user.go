package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       uint32    `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Password []byte    `json:"password"`
	Email    string    `json:"email"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

func (user *User) SetPassword(password string) error {
	// password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return nil
}

func (user *User) ComparePassword(password string) error {
	// Comparing the hashed user.password with the input string password
	return bcrypt.CompareHashAndPassword((user.Password), []byte(password))
}
