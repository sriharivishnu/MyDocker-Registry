package models

import (
	"fmt"
	"time"
)

type User struct {
	Id        string    `json:"id,omitempty" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}

func (user User) Validate() error {
	if len(user.Username) < 5 {
		return fmt.Errorf("username must be at least 5 characters in length")
	}

	if len(user.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters in length")
	}

	return nil
}
