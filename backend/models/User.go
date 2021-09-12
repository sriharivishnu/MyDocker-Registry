package models

import "time"

type User struct {
	Id        string    `json:"id,omitempty" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password,omitempty" db:"password"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}

type IUser interface {
	Create(email string, password string) error
	GetByEmail(email string) error
	GetById(id string) error
}
