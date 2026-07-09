package service

import (
	"product-api/internal/model"
	"product-api/internal/repository"
	"github.com/redis/go-redis/v9"
	"encoding/json"
	"errors"
	"context"
	"time"
	"fmt"
)

type ProductService struct {
	repo repository.ProductRepository
	redis *redis.Client
}

func NewProductService(repo repository.ProductRepository, redis *redis.Client) *ProductService {
	return &ProductService{repo: repo, redis: redis}
}

func (s *ProductService) GetAllProducts() ([]model.Product, error) {
	ctx := context.Background()
	chacheKey := "products_all"

	chaced, err := s.redis.Get(ctx, chacheKey).Result()
	if err == nil {
		var products []model.Product
		if err := json.Unmarshal([]byte(chaced), &products); err == nil {
			fmt.Println("Data dari cache Redis")
			return products, nil
		}
	}
	
	fmt.Println("Data dari database MySQL")
	products, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(products)
	if err != nil {
		s.redis.Set(ctx, chacheKey, data, 10*time.Second)
	}
	return products, nil
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

	result, err := s.repo.Create(p) 
	if err != nil {
		return model.Product{}, err
	}

	ctx := context.Background()
	if err := s.redis.Del(ctx, "products_all").Err(); err != nil {
		fmt.Println("Gagal menghapus cache Redis:", err)
	}

	return result, nil
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

	result, err := s.repo.Update(id, p)
	if err != nil {
		return model.Product{}, err
	}

	ctx := context.Background()
	s.redis.Del(ctx, "products_all")
	s.redis.Del(ctx, fmt.Sprintf("product_%d", id))

	return result, nil
}

func (s *ProductService) DeleteProduct(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	ctx := context.Background()
	s.redis.Del(ctx, "products_all")
	s.redis.Del(ctx, fmt.Sprintf("product_%d", id))

	return nil
}