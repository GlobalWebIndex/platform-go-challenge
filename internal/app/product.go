package app

import (
	"context"
	"fmt"
	"log"
	"ownify_api/internal/domain"
	"ownify_api/internal/dto"
	desc "ownify_api/pkg"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (m *MicroserviceServer) AddProduct(ctx context.Context, req *desc.AddProductRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Println(md.Get("test"))
	token, err := m.getUserIdFromToken(ctx)
	if err != nil {
		log.Println("user isn't authorized")
		return nil, err
	}
	_, err = m.tokenManager.ValidateFirebase(token)
	if err != nil {
		log.Println("user isn't authorized")
		return nil, err
	}

	// add product.
	product := dto.BriefProduct{
		ChainId:        int(req.ChainId),
		AssetId:        int64(req.AssetId),
		Owner:          req.Owner,
		Barcode:        req.Barcode,
		ItemName:       req.ItemName,
		BrandName:      req.BrandName,
		AdditionalData: req.AdditionalData,
		Location:       req.Location,
		IssueDate:      req.IssueDate,
	}

	err = m.productService.AddProduct(product, req.Net, true)
	if err != nil {
		return nil, err
	}
	return &desc.NetWorkResponse{
		Msg:     "Successfully Added",
		Success: true,
	}, nil
}

func (m *MicroserviceServer) AddProducts(ctx context.Context, req *desc.AddProductsRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	// md, _ := metadata.FromIncomingContext(ctx)
	// fmt.Println(md.Get("test"))
	// token, err := m.getUserIdFromToken(ctx)
	// if err != nil {
	// 	log.Println("user isn't authorized")
	// 	return nil, err
	// }
	// _, err = m.tokenManager.ValidateFirebase(token)
	// if err != nil {
	// 	log.Println("user isn't authorized")
	// 	return nil, err
	// }

	if req.Net != domain.TestNet || req.Net == domain.MainNet {
		return nil, fmt.Errorf("invalid network: %s", req.Net)
	}

	products := []dto.BriefProduct{}
	dupRemover := make(map[int64]int)
	for index, product := range req.Products {
		if _, ok := dupRemover[product.AssetId]; ok {
			return nil, fmt.Errorf("[ERR] include duplicated product information at %d", index)
		}

		dupRemover[product.AssetId] = 1
		product := dto.BriefProduct{
			ChainId:        int(req.ChainId),
			AssetId:        int64(product.AssetId),
			Owner:          product.Owner,
			Barcode:        product.Barcode,
			ItemName:       product.ItemName,
			BrandName:      product.BrandName,
			AdditionalData: product.AdditionalData,
			Location:       product.Location,
			IssueDate:      product.IssueDate,
		}

		if !product.Valid() {
			return nil, fmt.Errorf("[ERR] include invalid product information at %d", index)
		}
		products = append(products, product)
	}
	// add product.
	err := m.productService.AddProducts(products, req.Net, true)
	if err != nil {
		return nil, err
	}

	return &desc.NetWorkResponse{
		Msg:     "Successfully Added",
		Success: true,
	}, nil
}

func (m *MicroserviceServer) GetProduct(ctx context.Context, req *desc.SignInRequest) (*emptypb.Empty, error) {
	// md, _ := metadata.FromIncomingContext(ctx)
	// fmt.Println(md.Get("test"))
	// token, err := m.authService.SignIn(req.GetEmail(), req.GetPassword())
	// if err != nil {
	// 	return nil, err
	// }
	// err = grpc.SendHeader(ctx, metadata.New(map[string]string{
	// 	"Token": *token,
	// }))

	// if err != nil {
	// 	return nil, err
	// }
	return &emptypb.Empty{}, nil
}

func (m *MicroserviceServer) VerifyProduct(ctx context.Context, req *desc.VerifyAssetRequest) (*desc.NetWorkResponse, error) {
	product, err := m.productService.GetProduct(
		int(req.ChainId),
		req.AssetId,
		req.Net,
	)

	if err != nil {
		return nil, err
	}
	return BuildRes(product, "successfully verified", true)
}

func (m *MicroserviceServer) GetProducts(ctx context.Context, req *desc.GetProductsRequest) (*desc.NetWorkResponse, error) {
	products, err := m.productService.GetProducts(
		req.Net,
		int(req.Page),
		int(req.PerPage),
	)
	if err != nil {
		return nil, err
	}

	type data struct {
		Products []dto.BriefProduct `json:"products"`
	}

	return BuildRes(data{Products: products}, "there are your products", true)
}

func (m *MicroserviceServer) SearchProducts(ctx context.Context, req *desc.SearchProductsRequest) (*desc.NetWorkResponse, error) {

	filter := dto.BriefProduct{
		AssetId:   req.Filter.AssetId,
		Owner:     req.Filter.Owner,
		Barcode:   req.Filter.Barcode,
		ItemName:  req.Filter.ItemName,
		BrandName: req.Filter.BrandName,
		IssueDate: req.Filter.IssueDate,
	}
	products, err := m.productService.SearchProducts(filter, req.Net, req.Page, req.PerPage)

	if err != nil {
		return nil, err
	}

	type data struct {
		Products []dto.BriefProduct `json:"products"`
	}

	return BuildRes(data{Products: products}, "there are your products", true)
}
