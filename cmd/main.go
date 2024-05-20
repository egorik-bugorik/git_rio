package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"pgxrio/internal/database"
	"pgxrio/internal/inventory"
	"pgxrio/internal/postgres"
)

var (
	httpAddr = flag.String("http", "localhost:8080", "Addres for http server")
	grpcAddr = flag.String("grpc", "localhost:8082", "Addres for grpc server")
)

func main() {
	flag.Parse()

	connStr := "host=localhost port=5432 user=postgres password=1  database=pgxtutorial"

	ctx := context.Background()
	var logger = database.PGXStdLogger{}
	logLvl, err := database.LogLevelFromConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	pgPool, err := database.NewPgxPool(ctx, connStr, &logger, logLvl)
	if err != nil {
		panic(err)

	}

	err = pgPool.Ping(ctx)
	if err != nil {
		panic(err)

	}
	db := postgres.NewDB(pgPool)
	//createTest(err, db, ctx)
	//searchTest(err, db, ctx)
	//getTest(err, db, ctx)
	name := "sheep"
	upda := inventory.UpdateProductParams{
		ID:          "123",
		Name:        &name,
		Description: nil,
		Price:       nil,
	}
	err = db.UpdateProduct(ctx, upda)
	if err != nil {
		panic(err)
	}
	var _ = pgPool

}

func searchTest(err error, db postgres.DB, ctx context.Context) {
	sea := inventory.SearchProductParams{
		QueryString: "vd",
		MinPrice:    0,
		MaxPrice:    0,
		Pagination:  inventory.Pagination{},
	}
	product, err := db.SearchProduct(ctx, sea)
	fmt.Println("\n **********\n ", product)
}

func getTest(err error, db postgres.DB, ctx context.Context) {
	product, err := db.GetProduct(ctx, "123")
	if err != nil {
		panic(err)

	}
	fmt.Println("\n *******\n ", product)
}

func createTest(err error, db postgres.DB, ctx context.Context) {
	crea := inventory.CreateProductParams{
		ID:          "123231",
		Name:        "Ball",
		Description: "for fun",
		Price:       0,
	}
	err = db.CreateProduct(ctx, crea)
	if err != nil {
		panic(err)

	}
}
