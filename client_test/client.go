package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"pgxrio/internal/api"
)

func main() {

	conn, err := grpc.NewClient(":8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)

	}

	defer conn.Close()

	client := api.NewInventoryClient(conn)
	//err = createTest(err, client)
	prods, err := client.SearchProducts(context.Background(), &api.SearchProductsRequest{
		QueryString: "",
		MinPrice:    nil,
		MaxPrice:    nil,
		Page:        nil,
	})
	if err != nil {
		panic(err)

	}
	for _, item := range prods.Items {
		log.Println(":::   ", item.Name, "    :::")

	}

	log.Println("OK!")
}

func createTest(err error, client api.InventoryClient) error {
	_, err = client.CreateProduct(context.Background(), &api.CreateProductRequest{
		Id:          "000",
		Name:        "Metal cuup",
		Description: "For tea or coffe",
		Price:       999,
	})
	return err
}
