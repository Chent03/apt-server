package controllers

import (
	"net/http"

	"github.com/chent03/apt-server/context"
	"github.com/chent03/apt-server/models"
)

type Reviews struct {
	rs models.ReviewService
}

type ReviewForm struct {
	Review string `json:"review"`
}

func NewReviews(rs models.ReviewService) *Reviews {
	return &Reviews{
		rs: rs,
	}
}

func (re *Reviews) Create(w http.ResponseWriter, r *http.Request) {
	var form ReviewForm
	err := parseResponse(r, &form)
	if err != nil {
		RespondWithPayload(w, http.StatusBadRequest, &Payload{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}
	user := context.User(r.Context())
	review := models.Review{
		UserID: user.ID,
		Review: form.Review,
	}
	if err := re.rs.Create(&review); err != nil {
		RespondWithPayload(w, http.StatusBadRequest, &Payload{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}
	RespondWithPayload(w, http.StatusAccepted, &Payload{
		Success: true,
	})
}
