package main

import (
	"fmt"
	"log"
	"net"

	"github.com/syntropy-workshop/uplink-server/uplink_provider"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalln("Error: Failed to listen on port 9000: %w", err)
	}

	grpcServer := grpc.NewServer()

	s := uplink_provider.Server{}

	uplink_provider.RegisterUplinkServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Error: Failed to serve gRPC on port 9000: %w", err)
	}

	fmt.Println("Listening on port 9000")
}
