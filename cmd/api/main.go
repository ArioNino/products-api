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

	_ "product-api/docs"
)

func main() {
	repo := repository.NewProductRepository()
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)

	mux := router.NewRouter(h)

	fmt.Println("Server jalan di port 8081")
	fmt.Println("Swagger docs di http://localhost:8081/swagger/index.html")
	log.Fatal(http.ListenAndServe(":8081", mux))
}