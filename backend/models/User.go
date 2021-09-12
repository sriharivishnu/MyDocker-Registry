package models

import (
	"log"
	"time"

	"github.com/pkg/errors"
	db "github.com/sriharivishnu/shopify-challenge/services"
)

type User struct {
	Id        string    `json:"id,omitempty" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}

type IUser interface {
	Create(username string, password string) error
	GetByUsername(username string) error
	GetById(id string) error
}

func (user *User) Create() error {
	tx := db.DbConn.MustBegin()
	log.Printf("%s %s", user.Username, user.Password)
	err := tx.Get(user, "INSERT INTO user (username, password) VALUES (?, ?) returning *;", user.Username, user.Password)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "create user error")
	}
	tx.Commit()
	return nil
}

func (user *User) GetByUsername(username string) error {
	err := db.DbConn.Get(user, "select * from user where username = ?", username)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) GetById(id string) error {
	err := db.DbConn.Get(user, "select * from user where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
