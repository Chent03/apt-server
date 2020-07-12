package controllers

import (
	"net/http"

	"github.com/chent03/apt-server/models"
)

type Users struct {
	us *models.UserService
}

type SignupForm struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUsers(us *models.UserService) *Users {
	return &Users{
		us: us,
	}
}

func (u *Users) Register(w http.ResponseWriter, r *http.Request) {
	var form SignupForm

}
