package models

import (
	"errors"

	"github.com/chent03/apt-server/hash"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	FirstName    string `json:"firstName"; gorm:"not null";`
	LastName     string `json:"lastName"; gorm:"not null";`
	Email        string `json:"email"; gorm:"not null; unique_index";`
	Password     string `gorm:"-" json:"-"`
	PasswordHash string `gorm:"not null" json:"-"`
	Remember     string `gorm:"-" json:"-"`
	RememberHash string `gorm:"not null; unique_index" json:"-"`
	IsAdmin      bool   `json:"isAdmin";`
}

var (
	ErrNotFound        = errors.New("models: resource not found")
	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

var userPwPepper = "secret-random-string"

const hmacSecretKey = "secret-hmac-key"

type UserDB interface {
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}

type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

type userService struct {
	UserDB
}

// NewUserService opens a new db connection
func NewUserService(db *gorm.DB) UserService {
	ug := &userGorm{db}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)
	return &userService{
		UserDB: uv,
	}
}

func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash),
		[]byte(password+userPwPepper))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrInvalidPassword
	default:
		return nil, err
	}
}
