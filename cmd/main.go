package main

import (
	"context"
	"flag"
	"pgxrio/internal/database"
)

var (
	httpAddr = flag.String("http", "localhost:8080", "Addres for http server")
	grpcAddr = flag.String("grpc", "localhost:8082", "Addres for grpc server")
)

func main() {
	flag.Parse()

	connStr := "host=localhost port=5432 user=postgres password=1  database=pgxtutorial"

	ctx := context.Background()
	var logger tracer.Logger
	logLvl :=
	pgPool, err := database.NewPgxPool(ctx, connStr,logger,logLvl)

}
