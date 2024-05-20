package inventory

import "time"

type Review struct {
	ID          string
	ProductID   string
	ReviewerID  string
	Score       int
	Title       string
	Description string
	CreatedAt   time.Time
	ModifiedAt  time.Time
}
type CreateReviewParams struct {
	ID          string
	Name        string
	Description string
	Price       int
}

func (p *CreateReviewParams) validate() error {
	return nil
}

type SearchReviewParams struct {
	ID          string
	Name        string
	Description string
	Price       int
}

func (p *SearchReviewParams) validate() error {
	return nil
}

type UpdateReviewParams struct {
	ID          string
	Name        *string
	Description *string
	Price       *int
}

func (p *UpdateReviewParams) validate() error {
	return nil
}

type SearchReviewResponse struct {
	Items []*Review
	Count int
}
