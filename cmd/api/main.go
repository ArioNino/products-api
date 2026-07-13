// Package main Product API
//
// @title           Product API
// @version         1.0
// @description     REST API untuk CRUD produk menggunakan net/http dan ServeMux
// @host            localhost:8081
// @BasePath        /
package main

import (
	"fmt"
	"log"
	"log/slog"
	"context"
	"net/http"
	"os"
	_ "product-api/docs"
	"product-api/internal/database"
	"product-api/internal/handler"
	"product-api/internal/observability"
	"product-api/internal/repository"
	"product-api/internal/router"
	"product-api/internal/service"
)

func main() {
	dsn := "root:rootpassword@tcp(localhost:3306)/products_db"

	db, err := database.ConnectDB(dsn)
	if err != nil {
		log.Fatal(fmt.Errorf("gagal membuat koneksi MySQL: %w", err))
	}
	defer db.Close()

	redisClient, err := database.ConnectRedis("localhost:6379")
	if err != nil {
		log.Fatal(fmt.Errorf("gagal membuat koneksi Redis: %w", err))
	}
	defer redisClient.Close()

	logger, shutdownLogger, err := observability.InitLogger("product-api")
	if err != nil {
		log.Fatal(fmt.Errorf("gagal setup logger: %w", err))
	}
	defer shutdownLogger(context.Background())

	slog.SetDefault(logger)

	repo := repository.NewProductRepositoryMySQL(db)
	svc := service.NewProductService(repo, redisClient, logger)
	h := handler.NewProductHandler(svc)

	mux := router.NewRouter(h)

	if err := http.ListenAndServe(":8081", mux); err != nil {
		slog.Error("server gagal berjalan", "error", err)
		os.Exit(1)
	}
}
