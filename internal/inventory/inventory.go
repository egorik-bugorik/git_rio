package inventory

import (
	"context"
	"errors"
	"time"
)

type Product struct {
	ID          string
	Name        string
	Description string
	Price       int
	CreatedAt   time.Time
	ModifiedAt  time.Time
}
type CreateProductParams struct {
	ID          string
	Name        string
	Description string
	Price       int
}

func (p *CreateProductParams) validate() error {
	return nil
}

type SearchProductParams struct {
	QueryString string
	MinPrice    int
	MaxPrice    int
	Pagination  Pagination
}

func (p *SearchProductParams) validate() error {
	return nil
}

type UpdateProductParams struct {
	ID          string
	Name        *string
	Description *string
	Price       *int
}

func (p *UpdateProductParams) validate() error {
	return nil
}

type SearchProductResponse struct {
	Items []*Product
	Count int
}

func (s *Service) CreateProduct(ctx context.Context, params CreateProductParams) error {
	if err := params.validate(); err != nil {
		return err
	}
	return s.db.CreateProduct(ctx, params)
}

func (s *Service) SearchProduct(ctx context.Context, params SearchProductParams) (*SearchProductResponse, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	return s.db.SearchProduct(ctx, params)
}

func (s *Service) UpdateProduct(ctx context.Context, params UpdateProductParams) error {
	if err := params.validate(); err != nil {
		return err
	}
	return s.db.UpdateProduct(ctx, params)
}

func (s *Service) GetProduct(ctx context.Context, id string) (*Product, error) {
	if id == "" {
		return nil, errors.New("wrong id!!!")
	}
	return s.db.GetProduct(ctx, id)
}

func (s *Service) DeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("wrong id!!!")
	}
	return s.db.DeleteProduct(ctx, id)
}

func (s *Service) CreateReview(ctx context.Context, params CreateReviewParams) error {
	return s.db.CreateReview(ctx, params)
}

func (s *Service) SearchReview(ctx context.Context, params SearchReviewParams) (SearchReviewResponse, error) {
	return s.db.SearchReview(ctx, params)
}

func (s *Service) UpdateReview(ctx context.Context, params UpdateReviewParams) error {
	return s.db.UpdateReview(ctx, params)
}

func (s *Service) GetReview(ctx context.Context, id string) (Review, error) {
	return s.db.GetReview(ctx, id)
}

func (s *Service) DeleteReview(ctx context.Context, id string) error {
	return s.db.DeleteReview(ctx, id)
}
