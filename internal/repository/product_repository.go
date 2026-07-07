package repository

import (
	"errors"
	"sync"

	"product-api/internal/model"
)

type ProductRepository struct {
	mu       sync.Mutex
	products map[int]model.Product
	nextID   int
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		products: make(map[int]model.Product),
		nextID:   1,
	}
}

func (r *ProductRepository) Create(p model.Product) model.Product {
	r.mu.Lock()
	defer r.mu.Unlock()

	p.ID = r.nextID
	r.products[p.ID] = p
	r.nextID++
	return p
}

func (r *ProductRepository) GetAll() []model.Product {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]model.Product, 0, len(r.products))
	for _, p := range r.products {
		result = append(result, p)
	}
	return result
}

func (r *ProductRepository) GetByID(id int) (model.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, ok := r.products[id]
	if !ok {
		return model.Product{}, errors.New("produk tidak ditemukan")
	}
	return p, nil
}

func (r *ProductRepository) Update(id int, updated model.Product) (model.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.products[id]
	if !ok {
		return model.Product{}, errors.New("produk tidak ditemukan")
	}

	updated.ID = id
	r.products[id] = updated
	return updated, nil
}

func (r *ProductRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.products[id]
	if !ok {
		return errors.New("produk tidak ditemukan")
	}

	delete(r.products, id)
	return nil
}