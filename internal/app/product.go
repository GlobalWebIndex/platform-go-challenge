package app

import (
	"context"
	"fmt"
	"log"
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

	err = m.productService.AddProduct(product, req.Net)
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
