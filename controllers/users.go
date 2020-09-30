package controllers

import (
	"net/http"

	"github.com/chent03/apt-server/context"
	"github.com/chent03/apt-server/models"
	"github.com/chent03/apt-server/rand"
)

type Users struct {
	us models.UserService
}

type SignupForm struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Payload struct {
	Success      bool        `json:"success"`
	Data         interface{} `json:"data"`
	ErrorMessage string      `json:"message,omitempty"`
}

func NewUsers(us models.UserService) *Users {
	return &Users{
		us: us,
	}
}

func (u *Users) Register(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	err := parseResponse(r, &form)
	if err != nil {
		RespondWithPayload(w, http.StatusBadRequest, &Payload{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}
	user := models.User{
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Email:     form.Email,
		Password:  form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		RespondWithPayload(w, http.StatusBadRequest, &Payload{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}
	RespondWithPayload(w, http.StatusAccepted, &Payload{
		Success: true,
		Data:    user,
	})
}

func (u *Users) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	RespondWithPayload(w, http.StatusAccepted, &Payload{
		Success: true,
		Data:    user,
	})
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var form LoginForm
	err := parseResponse(r, &form)
	if err != nil {
		RespondWithPayload(w, http.StatusBadRequest, &Payload{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		RespondWithPayload(w, http.StatusBadRequest, &Payload{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		// switch err {
		// case models.ErrNotFound:
		// 	fmt.Println(w, "Invalid email address")
		// case models.ErrInvalidPassword:
		// 	fmt.Fprintln(w, "Invalid password provided")
		// default:
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// }
		return
	}
	err = u.signIn(w, user)
	if err != nil {
		RespondWithPayload(w, http.StatusBadRequest, &Payload{
			Success:      false,
			ErrorMessage: err.Error(),
		})
	}
	RespondWithPayload(w, http.StatusAccepted, &Payload{
		Success: true,
		Data:    user,
	})
}

func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)
	return nil
}
