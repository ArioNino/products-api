package service_test

import (
	"errors"
	"product-api/internal/model"
	// "product-api/internal/repository"

	// "product-api/internal/repository"
	"product-api/internal/repository/mocks"
	"product-api/internal/service"
	"reflect"
	"testing"

	// "vendor/golang.org/x/net/idna"

	"github.com/golang/mock/gomock"
	// "golang.org/x/tools/go/analysis/passes/defers"
)

func TestProductService_GetProductByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProductRepository(ctrl)

	expected := model.Product{
		ID: 1,
		Name: "Kopi Arabika",
		Price: 10000,
		Stock: 10,
	}

	tests := []struct {
		name    string
		id      int
		mockSetup func()
		want    model.Product
		wantErr bool
	}{
		{
			name:    "Success",
			id:      1,
			mockSetup: func() {
				mockRepo.EXPECT().GetByID(1).Return(expected, nil)
			},
			want:    expected,
			wantErr: false,
		},
		{
			name:    "Not Found",
			id:      999,
			mockSetup: func() {
				mockRepo.EXPECT().GetByID(999).Return(model.Product{}, errors.New("product not found"))
			},
			want:    model.Product{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			s := service.NewProductService(mockRepo)

			got, err := s.GetProductByID(tt.id)

			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error = %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %+v, want = %+v", got, tt.want)
			}
		})
	}
}

func TestProductService_GetAllProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProductRepository(ctrl)

	expected := model.Product{
		ID:    1,
		Name:  "Kopi Arabika",
		Price: 10000,
		Stock: 10,
	}

	tests := []struct {
		name      string
		mockSetup func()
		want      []model.Product
		wantErr   bool
	}{
		{
			name: "Success",
			mockSetup: func() {
				mockRepo.EXPECT().GetAll().Return([]model.Product{expected}, nil)
			},
			want:    []model.Product{expected},
			wantErr: false,
		},
		{
			name: "Error",
			mockSetup: func() {
				mockRepo.EXPECT().GetAll().Return(nil, errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			s := service.NewProductService(mockRepo)

			got, gotErr := s.GetAllProducts()

			if tt.wantErr {
				if gotErr == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if gotErr != nil {
				t.Fatalf("unexpected error = %v", gotErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %+v, want = %+v", got, tt.want)
			}
		})
	}
}

func TestProductService_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProductRepository(ctrl)

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
		name      string
		req       model.ProductCreateRequest
		mockSetup func()
		want      model.Product
		wantErr   bool
	}{
		{
			name: "Success",
			req:  validReq,
			mockSetup: func() {
				mockRepo.EXPECT().
					Create(gomock.Any()).
					Return(createdProduct, nil)
			},
			want:    createdProduct,
			wantErr: false,
		},
		{
			name:      "EmptyName",
			req:       model.ProductCreateRequest{
				Name: "", 
				Price: 10000, 
				Stock: 10},
			mockSetup: func() {}, 
			wantErr:   true,
		},
		{
			name:      "InvalidPrice",
			req:       model.ProductCreateRequest{
				Name: "Kopi", 
				Price: 0, 
				Stock: 10},
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name:      "InvalidStock",
			req:       model.ProductCreateRequest{
				Name: "Kopi", 
				Price: 10000, 
				Stock: 0},
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name: "RepoError",
			req:  validReq,
			mockSetup: func() {
				mockRepo.EXPECT().
					Create(gomock.Any()).
					Return(model.Product{}, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			s := service.NewProductService(mockRepo)

			got, gotErr := s.CreateProduct(tt.req)

			if tt.wantErr {
				if gotErr == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if gotErr != nil {
				t.Fatalf("unexpected error = %v", gotErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %+v, want = %+v", got, tt.want)
			}
		})
	}
}


func TestProductService_UpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProductRepository(ctrl)

	validReq := model.ProductUpdateRequest{
		Name: "kopi arabika",
		Price: 15000,
		Stock: 10,
	}

	updateProduct := model.Product{
		ID:		1,
		Name:	"Kopi arabika",
		Price:	15000,
		Stock: 	10,
	}

	tests := []struct {
		name      string
		id        int
		req       model.ProductUpdateRequest
		mockSetup func()
		want      model.Product
		wantErr   bool
	}{
		{
			name:	"Success",
			id:		1,
			req:	validReq,
			mockSetup: func ()  {
				mockRepo.EXPECT().
					Update(1, gomock.Any()).
					Return(updateProduct, nil)
			},
			want: updateProduct,
			wantErr : false,
		},
		{
			name:      "EmptyName",
			id:        1,
			req:       model.ProductUpdateRequest{
				Name: "",
				Price: 10000,
				Stock: 15},
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name:	"InvalidPrice",
			id:        1,
			req:       model.ProductUpdateRequest{
				Name: "Kopi Arabica",
				Price: -10000,
				Stock: 15},
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name:	"InvalidPrice",
			id:        1,
			req:       model.ProductUpdateRequest{
				Name: "Kopi Arabica",
				Price: 10000,
				Stock: -15},
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name: "RepoError",
			id:   999,
			req:  validReq,
			mockSetup: func() {
				mockRepo.EXPECT().
					Update(999, gomock.Any()).
					Return(model.Product{}, errors.New("produk tidak ditemukan"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			s := service.NewProductService(mockRepo)

			got, gotErr := s.UpdateProduct(tt.id, tt.req)

			if tt.wantErr {
				if gotErr == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if gotErr != nil {
				t.Fatalf("unexpected error = %v", gotErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %+v, want = %+v", got, tt.want)
			}
		})
	}
}


func TestProductService_DeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProductRepository(ctrl)

	tests := []struct {
		name      string
		id        int
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Success",
			id:   1,
			mockSetup: func() {
				mockRepo.EXPECT().
					Delete(1).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "RepoError",
			id:   999,
			mockSetup: func() {
				mockRepo.EXPECT().
					Delete(999).
					Return(errors.New("produk tidak ditemukan"))
			},
			wantErr: true,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			s := service.NewProductService(mockRepo)

			gotErr := s.DeleteProduct(tt.id)
			
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("unexpected error = %v, wantErr = %v", gotErr, tt.wantErr)
			}
		})
	}
}
