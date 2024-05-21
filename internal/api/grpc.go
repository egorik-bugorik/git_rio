package api

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"pgxrio/internal/inventory"
)

var _ InventoryServer = &InventoryGRPCServer{}

type InventoryGRPCServer struct {
	UnimplementedInventoryServer
	Service *inventory.Service
}

func (s *InventoryGRPCServer) SearchProducts(ctx context.Context, req *SearchProductsRequest) (*SearchProductsResponse, error) {
	params := inventory.SearchProductParams{
		QueryString: req.QueryString,
	}
	if req.MinPrice != nil {
		params.MinPrice = int(*req.MinPrice)
	}
	if req.MaxPrice != nil {
		params.MaxPrice = int(*req.MaxPrice)
	}
	page, pp := 1, 50
	if req.Page != nil {
		page = int(*req.Page)
	}
	params.Pagination = inventory.Pagination{
		Limit:  pp * page,
		Offset: pp * (page - 1),
	}
	products, err := s.Service.SearchProduct(ctx, params)
	if err != nil {
		return nil, grpcAPIError(err)
	}

	items := []*Product{}
	for _, p := range products.Items {
		items = append(items, &Product{
			Id:          p.ID,
			Price:       int64(p.Price),
			Name:        p.Name,
			Description: p.Description,
		})
	}
	return &SearchProductsResponse{
		Total: int32(products.Count),
		Items: items,
	}, nil
}
func (s *InventoryGRPCServer) CreateProduct(ctx context.Context, params *CreateProductRequest) (*CreateProductResponse, error) {

	if err := s.Service.CreateProduct(ctx, inventory.CreateProductParams{
		ID:          params.Id,
		Name:        params.Name,
		Description: params.Description,
		Price:       int(params.Price),
	}); err != nil {
		return nil, grpcAPIError(err)
	}

	return &CreateProductResponse{}, nil
}

func (s *InventoryGRPCServer) UpdateProduct(ctx context.Context, req *UpdateProductRequest) (*UpdateProductResponse, error) {
	params := inventory.UpdateProductParams{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
	}
	if req.Price != nil {
		price := int(*req.Price)
		params.Price = &price
	}
	if err := s.Service.UpdateProduct(ctx, params); err != nil {
		return nil, grpcAPIError(err)
	}
	return &UpdateProductResponse{}, nil
}
func (s *InventoryGRPCServer) DeleteProduct(ctx context.Context, params *DeleteProductRequest) (*DeleteProductResponse, error) {
	if err := s.Service.DeleteProduct(ctx, params.Id); err != nil {
		return nil, grpcAPIError(err)
	}
	return &DeleteProductResponse{}, nil
}
func (s *InventoryGRPCServer) GetProduct(ctx context.Context, params *GetProductRequest) (*GetProductResponse, error) {
	product, err := s.Service.GetProduct(ctx, params.Id)
	if err != nil {
		return nil, grpcAPIError(err)
	}
	if product == nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}
	return &GetProductResponse{
		Id:          product.ID,
		Price:       int64(product.Price),
		Name:        product.Name,
		Description: product.Description,
		CreatedAt:   product.CreatedAt.String(),
		ModifiedAt:  product.ModifiedAt.String(),
	}, nil
}
func grpcAPIError(err error) error {

	switch {
	case err == context.Canceled:
		return status.Error(codes.Canceled, err.Error())
	case err == context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, err.Error())
	case errors.As(err, &inventory.ValidationError{}):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return err

	}
}
