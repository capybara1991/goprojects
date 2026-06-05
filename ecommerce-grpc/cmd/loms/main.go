package main

import (
	"log"
	"net"
	"net/http"

	inventorypb "grpc_pr/gen/inventory"
	"grpc_pr/internal/httpx"
	"grpc_pr/internal/loms/controller"
	"grpc_pr/internal/loms/entity"
	lomshttp "grpc_pr/internal/loms/httpcontroller"
	"grpc_pr/internal/loms/repository"
	"grpc_pr/internal/loms/service"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	stocksRepo := repository.NewInMemoryStocksRepository()
	ordersRepo := repository.NewInMemoryOrderRepository()
	orderService := service.NewOrderService(ordersRepo, stocksRepo)

	products := map[uint32]entity.Product{
		1:  {Name: "Basic T-Shirt", Price: 1999},
		2:  {Name: "Hoodie", Price: 4999},
		3:  {Name: "Sneakers", Price: 8999},
		4:  {Name: "Backpack", Price: 3499},
		5:  {Name: "Cap", Price: 1299},
		6:  {Name: "Jeans", Price: 5999},
		7:  {Name: "Jacket", Price: 11999},
		8:  {Name: "Socks Pack", Price: 799},
		12: {Name: "T-shirt", Price: 27000},
		18: {Name: "Sneakers", Price: 49000},
		19: {Name: "Chair", Price: 11000},
	}

	lomsAPI := controller.NewAPI(orderService, stocksRepo, products)
	lomsHTTP := lomshttp.NewHandler(orderService, stocksRepo, products)

	httpMux := http.NewServeMux()
	lomsHTTP.Register(httpMux)
	go func() {
		log.Println("LOMS HTTP service started on :8081")
		if err := http.ListenAndServe(":8081", httpx.CORS(httpMux)); err != nil {
			log.Fatal(err)
		}
	}()

	grpcServer := grpc.NewServer()
	inventorypb.RegisterInventoryServiceServer(grpcServer, lomsAPI)

	log.Println("LOMS/Inventory service started on :50052")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
