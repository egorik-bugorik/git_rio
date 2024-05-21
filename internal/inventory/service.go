package inventory

import (
	"context"
)

type Service struct {
	db DB
}

func NewService(db DB) *Service {
	return &Service{db: db}
}

type ValidationError struct {
	s string
}

type Pagination struct {
	Limit  int
	Offset int
}
type DB interface {
	CreateProduct(ctx context.Context, params CreateProductParams) error
	SearchProduct(ctx context.Context, params SearchProductParams) (*SearchProductResponse, error)
	UpdateProduct(ctx context.Context, params UpdateProductParams) error
	GetProduct(ctx context.Context, id string) (*Product, error)
	DeleteProduct(ctx context.Context, id string) error
	CreateReview(ctx context.Context, params CreateReviewParams) error
	SearchReview(ctx context.Context, params SearchReviewParams) (SearchReviewResponse, error)
	UpdateReview(ctx context.Context, params UpdateReviewParams) error
	GetReview(ctx context.Context, id string) (Review, error)
	DeleteReview(ctx context.Context, id string) error
}
