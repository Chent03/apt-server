package middleware

import (
	"fmt"
	"net/http"

	"github.com/chent03/apt-server/context"
	"github.com/chent03/apt-server/controllers"
	"github.com/chent03/apt-server/models"
)

type RequireUser struct {
	models.UserService
}

func (mw *RequireUser) ApplFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			fmt.Println("No cookie")
			controllers.RespondWithPayload(w, http.StatusFound, &controllers.Payload{
				Success: false,
			})
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			controllers.RespondWithPayload(w, http.StatusFound, &controllers.Payload{
				Success: false,
			})
		}
		fmt.Println(user)
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next(w, r)
	})
}
