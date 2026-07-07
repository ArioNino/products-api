package router

import (
	"net/http"

	"product-api/internal/handler"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(productHandler *handler.ProductHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /products", productHandler.Create)
	mux.HandleFunc("GET /products", productHandler.GetList)
	mux.HandleFunc("GET /products/{id}", productHandler.GetDetail)
	mux.HandleFunc("PUT /products/{id}", productHandler.Update)
	mux.HandleFunc("DELETE /products/{id}", productHandler.Delete)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return mux
}