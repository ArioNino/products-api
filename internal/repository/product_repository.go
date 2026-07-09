package repository

import (
	"product-api/internal/model"
)

type ProductRepository interface {
	GetAll() ([]model.Product, error)
	Create(p model.Product) (model.Product, error)
	GetByID(id int) (model.Product, error)
	Update(id int, updated model.Product) (model.Product, error)
	Delete(id int) error
}