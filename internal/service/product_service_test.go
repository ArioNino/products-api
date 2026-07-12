package service_test

import (
	"encoding/json"
	"testing"
	cacheMock "product-api/internal/cache/mocks"
	"product-api/internal/model"
	repoMock "product-api/internal/repository/mocks"
	"product-api/internal/service"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
)

func TestProductService_GetAllProducts(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient)
		wantErr    bool
		wantLen    int
	}{
		{
			name: "CacheHit",
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {
				cachedProducts := []model.Product{
					{
						ID: 1, 
						Name: "Kopi Arabika",
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
			name: "Cache Miss Hit GetAll Products",
			setupMocks: func(mockRepo *repoMock.MockProductRepository, mockRedis *cacheMock.MockRedisClient) {
				mockRedis.EXPECT().
					Get(gomock.Any(), "products_all").
					Return(redis.NewStringResult("", redis.Nil))

				mockRepo.EXPECT().
					GetAll().
					Return([]model.Product{
						{
							ID: 1, 
							Name: "Kopi Arabika", 
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repoMock.NewMockProductRepository(ctrl)
			mockRedis := cacheMock.NewMockRedisClient(ctrl)

			tt.setupMocks(mockRepo, mockRedis)

			s := service.NewProductService(mockRepo, mockRedis)

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
		wantErr    bool
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
			name:       "return error when the name is empty ",
			req:        model.ProductCreateRequest{
				Name: "", 
				Price: 10000, 
				Stock: 10},
			setupMocks: func(mockRepo *repoMock.MockProductRepository, 
				mockRedis *cacheMock.MockRedisClient) {},
			wantErr:    true,
		},
		{
			name:       "return error when the price is invalid",
			req:        model.ProductCreateRequest{
				Name: "Kopi", 
				Price: 0, 
				Stock: 10},
			setupMocks: func(mockRepo *repoMock.MockProductRepository, 
				mockRedis *cacheMock.MockRedisClient) {},
			wantErr:    true,
		},
		{
			name:       "return error when the stock is invalid",
			req:        model.ProductCreateRequest{
				Name: "Kopi", 
				Price: 10000, 
				Stock: 0},
			setupMocks: func(mockRepo *repoMock.MockProductRepository, 
				mockRedis *cacheMock.MockRedisClient) {},
			wantErr:    true,
		},
		{
			name: "",
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

			s := service.NewProductService(mockRepo, mockRedis)

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


