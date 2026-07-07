package model

type Product struct {
	ID    int     `json:"id" example:"1"`
	Name  string  `json:"name" example:"Kopi Susu Ex"`
	Price float64 `json:"price" example:"15000"`
	Stock int     `json:"stock" example:"100"`
}

type ProductCreateRequest struct {
	Name  string  `json:"name" example:"Kopi Susu"`
	Price float64 `json:"price" example:"15000"`
	Stock int     `json:"stock" example:"100"`
}

type ProductUpdateRequest struct {
	Name  string  `json:"name" example:"Kopi Susu Gula Aren"`
	Price float64 `json:"price" example:"18000"`
	Stock int     `json:"stock" example:"80"`
}