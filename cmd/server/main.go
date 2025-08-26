package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "grpctasks/proto"
	"grpctasks/internal/server"
	"grpctasks/internal/store"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil { log.Fatal(err) }
	s := grpc.NewServer()
	st := store.NewMemory()
	srv := server.NewTaskServer(st)
	pb.RegisterTaskServiceServer(s, srv)
	log.Println("listening :50052")
	if err := s.Serve(lis); err != nil { log.Fatal(err) }
}
