package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Review struct {
	gorm.Model
	UserID uint   `gorm:"not_null";index`
	Review string `gorm:"not_null"`
}

type ReviewDB interface {
	ByID(id uint) (*Review, error)
	Create(review *Review) error
}

type ReviewService interface {
	ReviewDB
}

type reviewService struct {
	ReviewDB
}

type reviewValidator struct {
	ReviewDB
}

type reviewGorm struct {
	db *gorm.DB
}

type reviewValFn func(*Review) error

var (
	ErrUserIDRequired = errors.New("models: user ID is required")
	ErrReviewRequired = errors.New("models: review is required")
)

func NewReviewService(db *gorm.DB) ReviewService {
	return &reviewService{
		ReviewDB: &reviewValidator{
			ReviewDB: &reviewGorm{
				db: db,
			},
		},
	}
}

func (rg *reviewGorm) Create(review *Review) error {
	return rg.db.Create(review).Error
}

func (rg *reviewGorm) ByID(id uint) (*Review, error) {
	var review Review
	db := rg.db.Where("id = ?", id)
	err := first(db, &review)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func runReviewValFns(review *Review, fns ...reviewValFn) error {
	for _, fn := range fns {
		if err := fn(review); err != nil {
			return err
		}
	}
	return nil
}

func (rv *reviewValidator) userIDRequired(r *Review) error {
	if r.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (rv *reviewValidator) reviewRequired(r *Review) error {
	if r.Review == "" {
		return ErrReviewRequired
	}
	return nil
}

func (rv *reviewValidator) Create(review *Review) error {
	err := runReviewValFns(review,
		rv.userIDRequired,
		rv.reviewRequired,
	)
	if err != nil {
		return err
	}
	return rv.ReviewDB.Create(review)
}
