package repository

import (
	"errors"
	"sync"
	"product-api/internal/model"
)

type ProductRepositoryMemo struct {
	mu       sync.Mutex
	products map[int]model.Product
	nextID   int
}

func NewProductRepositoryMemo() *ProductRepositoryMemo {
	return &ProductRepositoryMemo{
		products: make(map[int]model.Product),
		nextID:   1,
	}
}

func (r *ProductRepositoryMemo) GetAll() ([]model.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]model.Product, 0, len(r.products))
	for _, p := range r.products {
		result = append(result, p)
	}
	return result, nil
}

func (r *ProductRepositoryMemo) Create(p model.Product) (model.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p.ID = r.nextID
	r.products[p.ID] = p
	r.nextID++
	return p, nil
}

func (r *ProductRepositoryMemo) GetByID(id int) (model.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, ok := r.products[id]
	if !ok {
		return model.Product{}, errors.New("produk tidak ditemukan")
	}
	return p, nil
}

func (r *ProductRepositoryMemo) Update(id int, updated model.Product) (model.Product, error) {
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

func (r *ProductRepositoryMemo) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.products[id]
	if !ok {
		return errors.New("produk tidak ditemukan")
	}

	delete(r.products, id)
	return nil
}