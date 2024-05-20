package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/henvic/pgtools"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"pgxrio/internal/database"
	"pgxrio/internal/inventory"
	"strings"
	"time"
)

var d inventory.DB = &DB{}

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

func (db *DB) CreateProduct(ctx context.Context, params inventory.CreateProductParams) error {

	q := `insert into products(id,name.description,price) values($1,$2,$3,$4)`
	switch _, err := db.conn(ctx).Exec(ctx, q, params.ID, params.Name, params.Description, params.Price); {
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		{
			return err
		}
	case err != nil:
		if pgxEr := db.productPgError(err); pgxEr != nil {
			return pgxEr
		}
		log.Print("coudln'tt create product on table ", err)
		err = fmt.Errorf("coudln'tt create product on table ")
		return err
	}
	return nil
}

func (db *DB) SearchProduct(ctx context.Context, params inventory.SearchProductParams) (*inventory.SearchProductResponse, error) {
	args := []any{"%" + params.QueryString + "%"}
	w := []string{"name LIKE $1"}

	if params.MinPrice != 0 {
		args = append(args, params.MinPrice)
		w = append(w, fmt.Sprintf("price>= $%d", len(args)))
	}
	if params.MaxPrice != 0 {
		args = append(args, params.MaxPrice)
		w = append(w, fmt.Sprintf("price<= $%d", len(args)))
	}

	where := strings.Join(w, " AND ")

	sqlTotal := fmt.Sprintf("SELECT COUNT(*) from products where  %s", where)

	resp := inventory.SearchProductResponse{
		Items: []*inventory.Product{},
	}
	switch err := db.conn(ctx).QueryRow(ctx, sqlTotal, args...).Scan(&resp.Count); {
	case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
		{
			return nil, err

		}
	case err != nil:
		log.Print("coudln'tt create product on table ", err)
		err = fmt.Errorf("coudln'tt create product on table ")

	}
	q := fmt.Sprintf(`SELECT * FROM product where $s order by id DESC`, where)
	if params.Pagination.Limit != 0 {
		args = append(args, params.Pagination.Limit)
		q += fmt.Sprintf("LIMIT $%d", len(args))
	}
	if params.Pagination.Offset != 0 {
		args = append(args, params.Pagination.Offset)
		q += fmt.Sprintf("OFFSET $%d", len(args))
	}
	rows, err := db.conn(ctx).Query(ctx, q, args...)
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil, err
	}
	var products []product
	if err == nil {
		products, err = pgx.CollectRows(rows, pgx.RowToStructByPos[product])
	}
	if err != nil {
		log.Print("coudln'tt search  products from table ", err)
		err = fmt.Errorf("coudln'tt search product from table ")
		return nil, err
	}
	for _, p := range products {
		var P *inventory.Product
		P = p.dto()
		resp.Items = append(resp.Items, P)

	}
	return &resp, err
}

var ErrorProductNotFound error = errors.New("Product not found!!!")

func (db *DB) UpdateProduct(ctx context.Context, params inventory.UpdateProductParams) error {
	q := `update products
SET 
    
name = COALESCE($1, "name"),
description = COALESCE($2, "description"),
price = COALESCE($3, "price"),
"modified_at" = now()
WHERE id = $4`
	ct, err := db.conn(ctx).Exec(ctx, q, params.Name, params.Description, params.Price, params.ID)
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	if pgxEr := db.productPgError(err); pgxEr != nil {
		return pgxEr
	}
	if err != nil {

		log.Print("coudln'tt update  product on table ", err)
		err = fmt.Errorf("coudln'tt update product on table ")
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrorProductNotFound
	}
	return nil

}

type product struct {
	ID          string
	Name        string
	Description string
	Price       int
	CreatedAt   time.Time
	ModifiedAt  time.Time
}

func (p *product) dto() *inventory.Product {
	return &inventory.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		ModifiedAt:  p.ModifiedAt,
	}
}

func (db *DB) GetProduct(ctx context.Context, id string) (*inventory.Product, error) {
	var p product
	q := fmt.Sprintf("SELECT %s from  product where id=$1", pgtools.Wildcard(p))

	rows, err := db.conn(ctx).Query(ctx, q, id)
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil, err
	}

	if err == nil {
		p, err = pgx.CollectOneRow(rows, pgx.RowToStructByPos[product])
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {

		log.Print("coudln'tt get  product from table ", err)
		err = fmt.Errorf("coudln'tt get product from table ")
		return nil, err
	}

	return p.dto(), nil

}

func (db *DB) DeleteProduct(ctx context.Context, id string) error {
	q := "DELETE  from  product where id=$1"

	_, err := db.conn(ctx).Exec(ctx, q, id)
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	if err != nil {

		log.Print("coudln'tt delete  product from table ", err)
		err = fmt.Errorf("coudln'tt delete product from table ")
		return err
	}

	return nil

}

func (db *DB) CreateReview(ctx context.Context, params inventory.CreateReviewParams) error {
	//TODO implement me
	panic("implement me")
}

func (db *DB) SearchReview(ctx context.Context, params inventory.SearchReviewParams) (inventory.SearchReviewResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) UpdateReview(ctx context.Context, params inventory.UpdateReviewParams) error {
	//TODO implement me
	panic("implement me")
}

func (db *DB) GetReview(ctx context.Context, id string) (inventory.Review, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) DeleteReview(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (db *DB) conn(ctx context.Context) database.PGXQuerier {

	if tx, ok := ctx.Value(txCtx{}).(pgx.Tx); tx != nil && ok {
		return tx
	} else if conn, ok := ctx.Value(connCtx{}).(*pgxpool.Conn); conn != nil && ok {
		return conn
	} else {
		return db.pool
	}

}

func (db *DB) productPgError(err error) error {

	var pgErr *pgconn.PgError
	if !errors.Is(err, pgErr) {
		return nil
	}
	switch pgErr.Code {
	case pgerrcode.UniqueViolation:
		return errors.New("Product already exists")
	case pgerrcode.CheckViolation:
		switch pgErr.ConstraintName {
		case "product_id_check":
			return errors.New("product's id is invalid")
		case "product_name_check":
			return errors.New("product's name is invalid")
		case "product_price_check":
			return errors.New("product's price is invalid")

		}

	}
	return nil

}
