package main

import (
	"context"
	"fmt"
	"log"
	"time"

	cartpb "grpc_pr/gen/cart"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cartClient := cartpb.NewCartServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := cartClient.AddItem(ctx, &cartpb.AddItemRequest{
		UserId: 1,
		Sku:    12,
		Count:  3,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response from cart:", resp.Message)
}
