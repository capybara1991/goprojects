package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"auth-service/internal/config"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	httptransport "auth-service/internal/transport/http"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	db, err := sqlx.Connect("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}
	defer db.Close()

	store := repository.NewRefreshStore(db.DB)
	authSvc := service.NewAuthService(
		cfg.JWTSecret,
		cfg.JWTExpireMinutes,
		cfg.TokenLengthBytes,
		cfg.RefreshExpireDays,
		cfg.BcryptCost,
		cfg.WebhookURL,
		store,
	)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/issue", httptransport.IssueHandler(authSvc))
	r.POST("/refresh", httptransport.RefreshHandler(authSvc))

	auth := r.Group("/")
	auth.Use(httptransport.AuthMiddleware(cfg.JWTSecret))
	auth.GET("/whoami", httptransport.WhoamiHandler())
	auth.POST("/logout", httptransport.LogoutHandler(authSvc))

	srv := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Server listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
