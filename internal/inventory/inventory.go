package inventory

import "time"

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
	ID          string
	Name        string
	Description string
	Price       int
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
