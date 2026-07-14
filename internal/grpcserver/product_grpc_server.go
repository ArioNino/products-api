package grpcserver

import (
	"context"
	"product-api/internal/model"
	"product-api/internal/service"
	pb "product-api/proto"
)


type ProductGRPCServer struct {
	pb.UnimplementedProductGRPCServiceServer
	service *service.ProductService
}

func NewProductGRPCServer (service *service.ProductService,) *ProductGRPCServer {
	return &ProductGRPCServer{service: service}
}

func toPbProduct(p model.Product) *pb.Product {
	return &pb.Product{
		Id: int32(p.ID),
		Name: p.Name,
		Price: p.Price,
		Stock: int32(p.Stock),
	}
}

func (s *ProductGRPCServer) GetAllProducts(ctx context.Context, req *pb.Empty) (*pb.ProductList, error) {
	products, err := s.service.GetAllProducts()
	if err != nil {
		return nil, err
	}

	var pbProducts []*pb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, toPbProduct(p))
	}

	return &pb.ProductList{Products: pbProducts}, nil
}

func (s *ProductGRPCServer) GetProductByID(ctx context.Context, req *pb.GetProductByIDRequest) (*pb.Product, error){
	p, err := s.service.GetProductByID(int(req.Id))
	if err != nil {
		return nil, err
	}

	return toPbProduct(p), nil
}

func (s *ProductGRPCServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error){
	created, err := s.service.CreateProduct(model.ProductCreateRequest{
		Name: req.Name,
		Price: req.Price,
		Stock: int(req.Stock),
	})
	if err != nil {
		return nil, err
	}

	return toPbProduct(created), nil
}

func (s *ProductGRPCServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error){
	updated, err := s.service.UpdateProduct(int(req.Id), model.ProductUpdateRequest{
		Name: req.Name,
		Price: req.Price,
		Stock: int(req.Stock),
	})
	if err != nil {
		return nil, err
	}

	return toPbProduct(updated), nil
}

func (s *ProductGRPCServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error){
	if err := s.service.DeleteProduct(int(req.Id)); err != nil{
		return nil, err
	}
	return &pb.DeleteProductResponse{Message: "product berhasil dihapus"}, nil
}
