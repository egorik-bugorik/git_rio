package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"pgxrio/internal/inventory"
)

var d inventory.DB = DB{}

type DB struct {
	pool *pgxpool.Pool
}

func (db *DB) TransactionContext(ctx context.Context) (context.Context, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, txCtx{}, tx), nil
}

func (db *DB) Commit(ctx context.Context) error {
	if tx, ok := ctx.Value(txCtx{}).(pgx.Tx); ok && tx != nil {
		return tx.Commit(ctx)
	}
	return errors.New("Context has no transactions!!!")
}

type txCtx struct {
}
type connCtx struct {
}

func (D DB) CreateProduct(ctx context.Context, params inventory.CreateProductParams) error {
	//TODO implement me
	panic("implement me")
}

func (D DB) SearchProduct(ctx context.Context, params inventory.SearchProductParams) (inventory.SearchProductResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (D DB) UpdateProduct(ctx context.Context, params inventory.UpdateProductParams) error {
	//TODO implement me
	panic("implement me")
}

func (D DB) GetProduct(ctx context.Context, id string) (inventory.Product, error) {
	//TODO implement me
	panic("implement me")
}

func (D DB) DeleteProduct(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (D DB) CreateReview(ctx context.Context, params inventory.CreateReviewParams) error {
	//TODO implement me
	panic("implement me")
}

func (D DB) SearchReview(ctx context.Context, params inventory.SearchReviewParams) (inventory.SearchReviewResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (D DB) UpdateReview(ctx context.Context, params inventory.UpdateReviewParams) error {
	//TODO implement me
	panic("implement me")
}

func (D DB) GetReview(ctx context.Context, id string) (inventory.Review, error) {
	//TODO implement me
	panic("implement me")
}

func (D DB) DeleteReview(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
