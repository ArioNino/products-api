package service_test

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	cacheMock "product-api/internal/cache/mocks"
	"product-api/internal/model"
	repoMock "product-api/internal/repository/mocks"
	"product-api/internal/service"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
)

var testLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestProductService_GetAllProducts(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient)
		wantErr    bool
		wantLen    int
	}{
		{
			name: "",
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {
				cachedProducts := []model.Product{
					{
						ID:    1,
						Name:  "Kopi Arabika",
						Price: 10000,
						Stock: 10},
				}
				cachedJSON, _ := json.Marshal(cachedProducts)

				mockRedis.EXPECT().
					Get(gomock.Any(), "products_all").
					Return(redis.NewStringResult(string(cachedJSON), nil))
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name: "return from mysql when cache miss",
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {
				mockRedis.EXPECT().
					Get(gomock.Any(), "products_all").
					Return(redis.NewStringResult("", redis.Nil))

				mockRepo.EXPECT().
					GetAll().
					Return([]model.Product{
						{
							ID:    1,
							Name:  "Kopi Arabika",
							Price: 10000,
							Stock: 10},
					}, nil)

				mockRedis.EXPECT().
					Set(gomock.Any(), "products_all", gomock.Any(), gomock.Any()).
					Return(redis.NewStatusResult("", nil)).
					AnyTimes()
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name: "return error when repository fails",
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {
				mockRedis.EXPECT().
					Get(gomock.Any(), "products_all").
					Return(redis.NewStringResult("", redis.Nil))

				mockRepo.EXPECT().
					GetAll().
					Return(nil, errors.New("koneksi database gagal"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repoMock.NewMockProductRepository(ctrl)
			mockRedis := cacheMock.NewMockRedisClient(ctrl)

			tt.setupMocks(mockRepo, mockRedis)

			s := service.NewProductService(mockRepo, mockRedis, testLogger)

			got, err := s.GetAllProducts()

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error = %v", err)
			}
			if len(got) != tt.wantLen {
				t.Errorf("expected %d products, got %d", tt.wantLen, len(got))
			}
		})
	}
}

func TestProductService_CreateProduct(t *testing.T) {
	validReq := model.ProductCreateRequest{
		Name:  "Kopi Arabika",
		Price: 10000,
		Stock: 10,
	}

	createdProduct := model.Product{
		ID:    1,
		Name:  "Kopi Arabika",
		Price: 10000,
		Stock: 10,
	}

	tests := []struct {
		name       string
		req        model.ProductCreateRequest
		setupMocks func(mockRepo *repoMock.MockProductRepository,
			mockRedis *cacheMock.MockRedisClient)
		wantErr bool
	}{
		{
			name: "create product success when input valid",
			req:  validReq,
			setupMocks: func(mockRepo *repoMock.MockProductRepository,
				mockRedis *cacheMock.MockRedisClient) {
				mockRepo.EXPECT().
					Create(gomock.Any()).
					Return(createdProduct, nil)

				mockRedis.EXPECT().
					Del(gomock.Any(), "products_all").
					Return(redis.NewIntResult(1, nil))
			},
			wantErr: false,
		},
		{
			name: "create product succes when cache fails",
			req: model.ProductCreateRequest{
				Name:  "Kopi Arabika",
				Price: 10000,
				Stock: 10,
			},
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {
				mockRepo.EXPECT().
					Create(gomock.Any()).
					Return(model.Product{ID: 1, Name: "Kopi Arabika", Price: 10000, Stock: 10}, nil)

				mockRedis.EXPECT().
					Del(gomock.Any(), "products_all").
					Return(redis.NewIntResult(0, errors.New("koneksi redis gagal")))
			},
			wantErr: false,
		},
		{
			name: "return error when the name is empty",
			req: model.ProductCreateRequest{
				Name:  "",
				Price: 10000,
				Stock: 10},
			setupMocks: func(mockRepo *repoMock.MockProductRepository,
				mockRedis *cacheMock.MockRedisClient) {
			},
			wantErr: true,
		},
		{
			name: "return error when the price is invalid",
			req: model.ProductCreateRequest{
				Name:  "Kopi",
				Price: 0,
				Stock: 10},
			setupMocks: func(mockRepo *repoMock.MockProductRepository,
				mockRedis *cacheMock.MockRedisClient) {
			},
			wantErr: true,
		},
		{
			name: "return error when the stock is invalid",
			req: model.ProductCreateRequest{
				Name:  "Kopi",
				Price: 10000,
				Stock: 0},
			setupMocks: func(mockRepo *repoMock.MockProductRepository,
				mockRedis *cacheMock.MockRedisClient) {
			},
			wantErr: true,
		},
		{
			name: "return error when repository fails",
			req:  validReq,
			setupMocks: func(mockRepo *repoMock.MockProductRepository,
				mockRedis *cacheMock.MockRedisClient) {
				mockRepo.EXPECT().
					Create(gomock.Any()).
					Return(model.Product{}, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repoMock.NewMockProductRepository(ctrl)
			mockRedis := cacheMock.NewMockRedisClient(ctrl)

			tt.setupMocks(mockRepo, mockRedis)

			s := service.NewProductService(mockRepo, mockRedis, testLogger)

			got, err := s.CreateProduct(tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error = %v", err)
			}
			if got.Name != createdProduct.Name {
				t.Errorf("expected name %s, got %s", createdProduct.Name, got.Name)
			}
		})
	}
}

func TestProductService_UpdateProduct(t *testing.T) {
	validReq := model.ProductUpdateRequest{
		Name:  "Kopi Arabika Update",
		Price: 12000,
		Stock: 15,
	}

	updatedProduct := model.Product{
		ID:    1,
		Name:  "Kopi Arabika Update",
		Price: 12000,
		Stock: 15,
	}

	tests := []struct {
		name       string
		id         int
		req        model.ProductUpdateRequest
		setupMocks func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient)
		wantErr    bool
	}{
		{
			name: "update product success when input valid",
			id:   1,
			req:  validReq,
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {
				mockRepo.EXPECT().
					Update(1, gomock.Any()).
					Return(updatedProduct, nil)

				mockRedis.EXPECT().
					Del(gomock.Any(), "products_all").
					Return(redis.NewIntResult(1, nil))

				mockRedis.EXPECT().
					Del(gomock.Any(), "product_1").
					Return(redis.NewIntResult(1, nil))
			},
			wantErr: false,
		},
		{
			name: "return error when name empty",
			id:   1,
			req: model.ProductUpdateRequest{
				Name:  "",
				Price: 12000,
				Stock: 15,
			},
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {},
			wantErr:    true,
		},
		{
			name: "return error when the price is zero",
			id:   1,
			req: model.ProductUpdateRequest{
				Name:  "Kopi Arabica",
				Price: 0,
				Stock: 15,
			},
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {},
			wantErr:    true,
		},
		{
			name: "return error when the stock is zero",
			id:   1,
			req: model.ProductUpdateRequest{
				Name:  "Kopi Arabica",
				Price: 12000,
				Stock: 0,
			},
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {},
			wantErr:    true,
		},
		{
			name: "return error when product not found",
			id:   999,
			req:  validReq,
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {
				mockRepo.EXPECT().
					Update(999, gomock.Any()).
					Return(model.Product{}, errors.New("produk tidak ditemukan"))
			},
			wantErr: true,
		},
		{
			name: "update product success",
			id:   1,
			req:  validReq,
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {
				mockRepo.EXPECT().
					Update(1, gomock.Any()).
					Return(updatedProduct, nil)

				mockRedis.EXPECT().
					Del(gomock.Any(), "products_all").
					Return(redis.NewIntResult(0, errors.New("redis connection error")))

				mockRedis.EXPECT().
					Del(gomock.Any(), "product_1").
					Return(redis.NewIntResult(0, errors.New("redis connection error")))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repoMock.NewMockProductRepository(ctrl)
			mockRedis := cacheMock.NewMockRedisClient(ctrl)

			tt.setupMocks(mockRepo, mockRedis)

			s := service.NewProductService(mockRepo, mockRedis, testLogger)

			got, err := s.UpdateProduct(tt.id, tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error = %v", err)
			}

			if got.Name != updatedProduct.Name {
				t.Errorf("expected name %s, got %s", updatedProduct.Name, got.Name)
			}
		})
	}
}

func TestProductService_DeleteProduct(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		setupMocks func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient)
		wantErr    bool
	}{
		{
			name: "delete product success when product exist",
			id:   1,
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {
				mockRepo.EXPECT().
					Delete(1).
					Return(nil)

				mockRedis.EXPECT().
					Del(gomock.Any(), "products_all").
					Return(redis.NewIntResult(1, nil))

				mockRedis.EXPECT().
					Del(gomock.Any(), "product_1").
					Return(redis.NewIntResult(1, nil))
			},
			wantErr: false,
		},
		{
			name: "return error when product not found",
			id:   999,
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {
				mockRepo.EXPECT().
					Delete(999).
					Return(errors.New("produk tidak ditemukan"))
			},
			wantErr: true,
		},
		{
			name: "delete product success when cache invalidation fails",
			id:   1,
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {
				mockRepo.EXPECT().
					Delete(1).
					Return(nil)

				mockRedis.EXPECT().
					Del(gomock.Any(), "products_all").
					Return(redis.NewIntResult(0, errors.New("redis connection error")))

				mockRedis.EXPECT().
					Del(gomock.Any(), "product_1").
					Return(redis.NewIntResult(0, errors.New("redis connection error")))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repoMock.NewMockProductRepository(ctrl)
			mockRedis := cacheMock.NewMockRedisClient(ctrl)

			tt.setupMocks(mockRepo, mockRedis)

			s := service.NewProductService(mockRepo, mockRedis, testLogger)

			err := s.DeleteProduct(tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error = %v, wanErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductService_GetProductByID(t *testing.T) {
	cachedProduct := model.Product{
		ID:    1,
		Name:  "Kopi Arabika",
		Price: 10000,
		Stock: 10,
	}

	tests := []struct {
		name       string
		id         int
		setupMocks func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient)
		wantErr    bool
	}{
		{
			name: "return product from cache when cache hit",
			id:   1,
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {
				cachedJSON, _ := json.Marshal(cachedProduct)

				mockRedis.EXPECT().
					Get(gomock.Any(), "product_1").
					Return(redis.NewStringResult(string(cachedJSON), nil))
			},
			wantErr: false,
		},
		{
			name: "return product from mysql when cache miss",
			id:   1,
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {
				mockRedis.EXPECT().
					Get(gomock.Any(), "product_1").
					Return(redis.NewStringResult("", redis.Nil))

				mockRepo.EXPECT().
					GetByID(1).
					Return(cachedProduct, nil)

				mockRedis.EXPECT().
					Set(gomock.Any(), "product_1", gomock.Any(), gomock.Any()).
					Return(redis.NewStatusResult("", nil)).
					AnyTimes()
			},
			wantErr: false,
		},
		{
			name: "return error when cache miss and product not found",
			id:   999,
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {
				mockRedis.EXPECT().
					Get(gomock.Any(), "product_999").
					Return(redis.NewStringResult("", redis.Nil))

				mockRepo.EXPECT().
					GetByID(999).
					Return(model.Product{}, errors.New("produk dengan id 999 tidak ditemukan"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repoMock.NewMockProductRepository(ctrl)
			mockRedis := cacheMock.NewMockRedisClient(ctrl)

			tt.setupMocks(mockRepo, mockRedis)

			s := service.NewProductService(mockRepo, mockRedis, testLogger)

			got, err := s.GetProductByID(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error = %v", err)
			}
			if got.Name != cachedProduct.Name {
				t.Errorf("expected name %s, got %s", cachedProduct.Name, got.Name)
			}
		})
	}
}
