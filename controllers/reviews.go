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
	Title   string `json:"title"`
	Review  string `json:"review"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address string `json:"address"`
	Rating  uint   `json:"rating"`
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
		UserID:  user.ID,
		Title:   form.Title,
		Review:  form.Review,
		Address: form.Address,
		Name:    form.Name,
		Type:    form.Type,
		Rating:  form.Rating,
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

func (re *Reviews) GetReviews(w http.ResponseWriter, r *http.Request) {
	reviews, err := re.rs.GetReviews()
	if err != nil {
		RespondWithPayload(w, http.StatusBadRequest, &Payload{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}
	RespondWithPayload(w, http.StatusAccepted, &Payload{
		Success: true,
		Data:    reviews,
	})
}

func (re *Reviews) GetByUserID(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	reviews, err := re.rs.ByUser(user.ID)
	if err != nil {
		RespondWithPayload(w, http.StatusBadRequest, &Payload{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}
	RespondWithPayload(w, http.StatusAccepted, &Payload{
		Success: true,
		Data:    reviews,
	})
}
