package main

import (
	"log"
	"net"
	"net/http"

	cartpb "grpc_pr/gen/cart"
	inventorypb "grpc_pr/gen/inventory"
	"grpc_pr/internal/cart/adapter"
	"grpc_pr/internal/cart/controller"
	carthttp "grpc_pr/internal/cart/httpcontroller"
	"grpc_pr/internal/cart/repository"
	"grpc_pr/internal/cart/service"
	"grpc_pr/internal/httpx"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	inventoryConn, err := grpc.Dial(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer inventoryConn.Close()

	inventoryClient := inventorypb.NewInventoryServiceClient(inventoryConn)
	productClient := adapter.NewProductClient(inventoryClient)
	lomsClient := adapter.NewLOMSClient(inventoryClient)
	cartRepository := repository.NewInMemoryCartRepository()
	itemService := service.NewItemService(cartRepository, productClient, lomsClient)
	cartAPI := controller.NewAPI(itemService)
	cartHTTP := carthttp.NewHandler(itemService)

	httpMux := http.NewServeMux()
	cartHTTP.Register(httpMux)
	go func() {
		log.Println("Cart HTTP service started on :8080")
		if err := http.ListenAndServe(":8080", httpx.CORS(httpMux)); err != nil {
			log.Fatal(err)
		}
	}()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	cartpb.RegisterCartServiceServer(grpcServer, cartAPI)

	log.Println("Cart service started on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
