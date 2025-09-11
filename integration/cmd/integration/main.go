package main

import (
	"log"
	"net"

	"integration/internal/config"
	"integration/internal/controller"
	"integration/internal/repository"
	"integration/internal/service"
	pb "integration/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	repo := repository.NewExternalSystem(cfg.ExternalURL)
	srv := service.NewPaymentService(repo)
	grpcServer := grpc.NewServer()

	pb.RegisterPaymentIntegrationServer(grpcServer, controller.NewPaymentController(srv))
	reflection.Register(grpcServer)

	lis, _ := net.Listen("tcp", ":"+cfg.GRPCPort)
	log.Println("🚀 Integration service started on", cfg.GRPCPort)
	grpcServer.Serve(lis)
}
