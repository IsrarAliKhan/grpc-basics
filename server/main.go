package main

import (
	"item/pb"
	"item/server/config"
	"item/server/middleware"
	"item/server/services"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	// create new server
	s := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.Unary()),
		grpc.StreamInterceptor(middleware.Stream()),
	)

	// register services
	pb.RegisterItemsServer(s, &services.ItemServer{})
	pb.RegisterAuthServer(s, &services.AuthServer{})

	// create listener
	l, e := net.Listen("tcp", config.RPCPort)
	if e != nil {
		log.Fatalf("Failed to listen: %v", e)
	}

	// start listening
	log.Printf("Server listening at %v", l.Addr())
	if e := s.Serve(l); e != nil {
		log.Fatalf("Failed to serve: %v", e)
	}
}
