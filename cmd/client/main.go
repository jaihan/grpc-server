package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"gitlab.com/techschool/pcbook/pb"
	"gitlab.com/techschool/pcbook/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createLaptop(laptopClient pb.LaptopServiceClient) {
	laptop := sample.NewLaptop()
	laptop.Id = ""

	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	// set timout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := laptopClient.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {

			log.Print("laptop already exists")
		} else {
			log.Fatal("cannot create laptop: ", err)
		}
		return
	}
	log.Printf("created laptop with id: %s ", res.Id)

}

func searchLaptop(laptopClient pb.LaptopServiceClient, filter *pb.Filter) {
	log.Print("search filter ", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("cannot search laptop:", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("cannot receive laptop:", err)
		}
		laptopt := res.GetLaptop()
		log.Print("- found :", laptopt.GetId())
		log.Print(" + brand:", laptopt.GetBrand())
		log.Print(" +name:", laptopt.GetName())
	}

}

func main() {
	serverAddress := flag.String("address", "", "server port")
	flag.Parse()
	log.Printf("dial server %d", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	laptopClient := pb.NewLaptopServiceClient(conn)

	for i := 0; i < 10; i++ {
		createLaptop(laptopClient)
	}

	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}

	searchLaptop(laptopClient, filter)

	// laptop := sample.NewLaptop()
	// req := &pb.CreateLaptopRequest{
	// 	Laptop: laptop,
	// }

	// res, err := laptopClient.CreateLaptop(context.Background(), req)

	// if err != nil {
	// 	st, ok := status.FromError(err)
	// 	if ok && st.Code() == codes.AlreadyExists {
	// 		log.Printf("laptop already exists")
	// 	} else {
	// 		log.Printf("laptop already exists")
	// 	}
	// 	return
	// }
	// log.Printf("created laptop with id: %s", res.Id)
}
