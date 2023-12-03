package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"gitlab.com/techschool/pcbook/pb"
	"gitlab.com/techschool/pcbook/service"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "server port")
	flag.Parse()
	log.Printf("start server on port %d", port)

	laptopStore := service.NewInMemoryLaptopStore()
	laptopServer := service.NewLaptopServer(laptopStore)
	grpcServer := grpc.NewServer()

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
