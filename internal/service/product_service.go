package service

import (
	"log/slog"
	"product-api/internal/model"
	"product-api/internal/repository"
	// "github.com/redis/go-redis/v9"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"product-api/internal/cache"
	"time"
)

type ProductService struct {
	repo repository.ProductRepository
	redis cache.RedisClient
	logger *slog.Logger
}

func NewProductService(repo repository.ProductRepository, redis cache.RedisClient, logger *slog.Logger) *ProductService {
	return &ProductService{repo: repo, redis: redis, logger: logger}
}

func (s *ProductService) GetAllProducts() ([]model.Product, error) {
	ctx := context.Background()
	cacheKey := "products_all"

	chaced, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var products []model.Product
		if err := json.Unmarshal([]byte(chaced), &products); err == nil {
			s.logger.Info("Data diambil dari cache", "source", "redis", "key", cacheKey)
			return products, nil
		}
	}

	s.logger.Info("Data diambil dari Mysql", "source", "mysql")
	products, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(products)
	if err == nil {
		s.redis.Set(ctx, cacheKey, data, 5*time.Minute)
	}
	return products, nil
}

func (s *ProductService) CreateProduct(req model.ProductCreateRequest) (model.Product, error) {
	if req.Name == "" {
		s.logger.Warn("validasi gagal: nama kosong", "req", req)
		return model.Product{}, errors.New("nama produk tidak boleh kosong")
	}
	if req.Price <= 0 {
		s.logger.Warn("validasi gagal: harga tidak valid", "req", req.Price)
		return model.Product{}, errors.New("harga tidak boleh negatif")
	}
	if req.Stock <= 0 {
		s.logger.Warn("validasi gagal: stock tidak valid", "req", req.Stock)
		return model.Product{}, errors.New("stok tidak boleh negatif")
	}

	p := model.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}

	result, err := s.repo.Create(p) 
	if err != nil {
		s.logger.Error("gagal membuat produk", "error", err)
		return model.Product{}, err
	}

	ctx := context.Background()
	if err := s.redis.Del(ctx, "products_all").Err(); err != nil {
		s.logger.Warn("gagal menghapus cache redis", "error", err)
	}

	s.logger.Info("produk berhasil dibuat", "product_id", result.ID, "name", result.Name)
	return result, nil
}

func (s *ProductService) GetProductByID(id int) (model.Product, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product_%d", id)

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var product model.Product
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			s.logger.Info("Data diambil dari cache", "source", "redis", "key", cacheKey)
			return product, nil
		}
	}

	s.logger.Info("Data diambil dari Mysql", "source", "mysql", "key", cacheKey)
	product, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Warn("produk tidak ditemukan", "product_id", id, "error", err)
		return model.Product{}, err
	}

	data, err := json.Marshal(product)
	if err == nil {
		s.redis.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return product, nil
}

func (s *ProductService) UpdateProduct(id int, req model.ProductUpdateRequest) (model.Product, error) {
	if req.Name == "" {
		s.logger.Warn("validasi gagal: nama kosong", "product_id", id)
		return model.Product{}, errors.New("nama produk tidak boleh kosong")
	}
	if req.Price <= 0 {
		s.logger.Warn("validasi gagal: harga tidak valid", "req", id, "price", req.Price)
		return model.Product{}, errors.New("harga tidak boleh negatif")
	}
	if req.Stock <= 0 {
		s.logger.Warn("validasi gagal: stock tidak valid", "req", id, "stock", req.Stock)
		return model.Product{}, errors.New("stok tidak boleh negatif")
	}

	p := model.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}

	result, err := s.repo.Update(id, p)
	if err != nil {
		s.logger.Error("gagal update produk", "product_id", id, "error", err)
		return model.Product{}, err
	}

	ctx := context.Background()
	if err := s.redis.Del(ctx, "products_all").Err(); err != nil {
		s.logger.Warn("gagal menghapus cache products_all", "error", err)
	}
	if err := s.redis.Del(ctx, fmt.Sprintf("product_%d", id)).Err(); err != nil {
		s.logger.Warn("gagal mengahpus cache produk", "product_id", id, "error", err)
	}

	s.logger.Info("produk berhasil diupdate", "product_id", id, "name", result.Name)
	return result, nil
}

func (s *ProductService) DeleteProduct(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	ctx := context.Background()
	if err := s.redis.Del(ctx, "products_all").Err(); err != nil {
		s.logger.Warn("gagal menghapus cache produk_all", "error", err)
	}
	if err := s.redis.Del(ctx, fmt.Sprintf("product_%d", id)).Err(); err != nil {
		s.logger.Warn("gagal mengahpus cache produk", "product_id", id, "error", err)
	}

	s.logger.Info("produk berhasil dihapus", "product_id", id)
	return nil
}