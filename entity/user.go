package entity

import (
	logger "go_rest_api_with_mysql/pkg/log"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	DEFAULTCOST int = 10
)

var log *zap.SugaredLogger = logger.GetLogger().Sugar()

// user data
type User struct {
	ID        ID
	Email     string
	Password  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// create a new user
func NewUser(userData *User) (*User, error) {
	user := &User{
		ID:        NewID(),
		Email:     userData.Email,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	pwd, err := generatePassword(userData.Password)
	if err != nil {
		log.Errorf("Error generating password. %v", err)
		return nil, err
	}
	user.Password = pwd
	err = user.Validate()
	if err != nil {
		log.Errorf("Invalid request entity data. %v", err)
		return nil, ErrInvalidEntity
	}
	return user, nil
}

// generate password
func generatePassword(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), DEFAULTCOST)
	if err != nil {
		log.Errorf("Error hashing password. %v", err)
		return "", err
	}
	return string(hash), err
}

// validate data
func (user *User) Validate() error {
	if user.Email == "" || user.FirstName == "" || user.LastName == "" {
		log.Errorf("Invalid entity")
		return ErrInvalidEntity
	}
	return nil
}
