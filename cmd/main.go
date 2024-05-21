package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"pgxrio/internal/api"
	"pgxrio/internal/database"
	"pgxrio/internal/inventory"
	"pgxrio/internal/postgres"
	"syscall"
	"time"
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
	pool, err := database.NewPgxPool(ctx, connStr, &logger, logLvl)
	s := &api.Server{
		HttpAddress: *httpAddr,
		GrpcAddress: *grpcAddr,
		Service:     inventory.NewService(postgres.NewDB(pool)),
	}

	errorChannel := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go func() {
		errorChannel <- s.Run(context.Background())

	}()

	select {
	case err = <-errorChannel:
	case <-ctx.Done():
		ctxHalt, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		s.Shutdown(ctxHalt)
		stop()
		err = <-errorChannel
		println("exit")
	}
	if err != nil {
		log.Fatal(err)
	}
}
