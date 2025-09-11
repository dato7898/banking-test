package main

import (
	"auth/internal/config"
	"auth/internal/controller"
	"auth/internal/repository"
	"auth/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg := config.Load()

	db, err := sqlx.Connect("pgx", cfg.PostgresDSN)
	if err != nil {
		log.Fatal("DB connect error:", err)
	}

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDb, cfg.SessionTTL)

	tokenService := service.NewTokenService(cfg.JWTSecret)
	authService := service.NewAuthService(userRepo, sessionRepo, tokenService)

	authHandler := controller.NewAuthHandler(authService)

	r := gin.Default()
	r.POST("/login", authHandler.Login)
	r.GET("/validate", authHandler.ValidateToken)

	log.Println("🚀 Auth Service started on", cfg.Port)
	r.Run(":" + cfg.Port)
}
