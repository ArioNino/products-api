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
	"net/http"
	"os"
	"product-api/internal/database"
	"product-api/internal/handler"
	"product-api/internal/repository"
	"product-api/internal/router"
	"product-api/internal/service"

	_ "product-api/docs"
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

	repo := repository.NewProductRepositoryMySQL(db)
	svc := service.NewProductService(repo, redisClient)
	h := handler.NewProductHandler(svc)

	mux := router.NewRouter(h)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	if err := http.ListenAndServe(":8081", mux); err != nil {
		slog.Error("server gagal berjalan", "error", err)
		os.Exit(1)
	}
}
