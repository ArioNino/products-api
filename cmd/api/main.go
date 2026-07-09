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
	"net/http"

	"product-api/internal/handler"
	"product-api/internal/repository"
	"product-api/internal/router"
	"product-api/internal/service"
	"product-api/internal/database"

	_ "product-api/docs"
)

func main() {
	dsn := "root:rootpassword@tcp(localhost:3306)/products_db"
	
	// MySQL
	db, err := database.ConnectDB(dsn)
	if err != nil {
		log.Fatal(fmt.Errorf("gagal membuat koneksi MySQL: %w", err))
	}
	defer db.Close()

	// Redis
	redisClient, err := database.ConnectRedis("localhost:6379")
	if err != nil {
		log.Fatal(fmt.Errorf("gagal membuat koneksi Redis: %w", err))
	}
	defer redisClient.Close()

	repo := repository.NewProductRepositoryMySQL(db)
	svc := service.NewProductService(repo, redisClient)
	h := handler.NewProductHandler(svc)

	mux := router.NewRouter(h)

	fmt.Println("Server jalan di port 8081")
	fmt.Println("Swagger docs di http://localhost:8081/swagger/index.html")
	log.Fatal(http.ListenAndServe(":8081", mux))
}