package main

import (
	"account/internal/config"
	"account/internal/controller"
	"account/internal/repository"
	"account/internal/service"
	"log"
	"net"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/jackc/pgx/v5/stdlib"

	pb "account/proto"
)

func main() {
	cfg := config.Load()

	db, err := sqlx.Connect("pgx", cfg.PostgresDSN)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	accountRepo := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepo)

	grpcServer := grpc.NewServer()
	pb.RegisterAccountClientServer(grpcServer, controller.NewAccountController(accountService))
	reflection.Register(grpcServer)

	lis, _ := net.Listen("tcp", ":"+cfg.GRPCPort)
	log.Println("🚀 Account service started on", cfg.GRPCPort)
	grpcServer.Serve(lis)
}
