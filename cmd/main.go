package main

import (
	"go_grpc_yt/cmd/config"
	"go_grpc_yt/cmd/services"
	productPB "go_grpc_yt/pb/product"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listened %v", err.Error())
	}

	grpcServer := grpc.NewServer()
	db := config.ConnectDB()

	// kita registerasika dulu product grpcnya ke dalam server
	prodService := services.ProductService{DB: db}
	productPB.RegisterProductServiceServer(grpcServer, &prodService)

	log.Printf("Server running at port:%v", listen.Addr())
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve %v ", err)
	}
}
