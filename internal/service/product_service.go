package service

import (
	"errors"

	"product-api/internal/model"
	"product-api/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(req model.ProductCreateRequest) (model.Product, error) {
	if req.Name == "" {
		return model.Product{}, errors.New("nama produk tidak boleh kosong")
	}
	if req.Price <= 0 {
		return model.Product{}, errors.New("harga tidak boleh negatif")
	}
	if req.Stock <= 0 {
		return model.Product{}, errors.New("stok tidak boleh negatif")
	}

	p := model.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}

	return s.repo.Create(p), nil
}

func (s *ProductService) GetAllProducts() []model.Product {
	return s.repo.GetAll()
}

func (s *ProductService) GetProductByID(id int) (model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) UpdateProduct(id int, req model.ProductUpdateRequest) (model.Product, error) {
	if req.Name == "" {
		return model.Product{}, errors.New("nama produk tidak boleh kosong")
	}
	if req.Price <= 0 {
		return model.Product{}, errors.New("harga tidak boleh negatif")
	}
	if req.Stock <= 0 {
		return model.Product{}, errors.New("stok tidak boleh negatif")
	}

	p := model.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}

	return s.repo.Update(id, p)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.Delete(id)
}