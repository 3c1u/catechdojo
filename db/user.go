package db

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NOTE: 一般的にはユーザーとセッションを分けるが，一つのユーザーに対して一つのトークンしか結びつかないため
// ここでは考えないことにする．

// User describes a model of a user.
type User struct {
	gorm.Model
	Name  string
	Token string
}

// UserCreate creates a user.
func UserCreate(name string) (user *User, err error) {
	if db == nil {
		err = fmt.Errorf("no database connection")
		return
	}

	// create a random token
	token, err := uuid.NewRandom()
	if err != nil {
		return
	}

	user = &User{
		Name:  name,
		Token: token.String(),
	}

	err = db.Create(user).Error
	if err != nil {
		user = nil
	}

	return
}

// UserGet finds a user with a given APi token.
func UserGet(token string) (user *User, err error) {
	if db == nil {
		err = fmt.Errorf("no database connection")
		return
	}

	user = &User{}
	if err = db.First(&user, "token = ?", token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		user = nil
		return
	}

	return
}

// UserUpdate updates a user info.
func UserUpdate(user *User) error {
	if db == nil {
		return fmt.Errorf("no database connection")
	}

	return db.Where("id = ?", user.ID).Save(&user).Error
}
