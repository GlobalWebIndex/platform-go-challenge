package app

import (
	"context"

	desc "ownify_api/pkg"
)

func (m *MicroserviceServer) CreateSubscription(ctx context.Context, req *desc.CreateSubscriptionRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	_, err := m.TokenInterceptor(ctx)
	if err != nil {
		return &desc.NetWorkResponse{}, err
	}

	subscriptionId, err := m.paymentService.CreateSubscription(req.CustomerId, req.PriceId)
	type CreateSubscriptionRes struct {
		Created        bool   `json:"created"`
		SubscriptionId string `json:"subscription_id"`
	}
	res := CreateSubscriptionRes{Created: true, SubscriptionId: *subscriptionId}
	if err != nil {
		res.Created = false
	}
	return BuildRes(res, "Successfully created!", true)

}

func (m *MicroserviceServer) UpdateSubscription(ctx context.Context, req *desc.UpdateSubscriptionRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	_, err := m.TokenInterceptor(ctx)
	if err != nil {
		return &desc.NetWorkResponse{}, err
	}
	err = m.paymentService.UpdateSubscription(req.OldPriceId, req.NewPriceId)

	type CreateSubscriptionRes struct {
		Created        bool   `json:"created"`
		SubscriptionId string `json:"subscription_id"`
	}
	res := CreateSubscriptionRes{Created: true}
	if err != nil {
		res.Created = false
	}
	return BuildRes(res, "Successfully created!", true)

}

func (m *MicroserviceServer) CancelSubscription(ctx context.Context, req *desc.CancelSubscriptionRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	_, err := m.TokenInterceptor(ctx)
	if err != nil {
		return &desc.NetWorkResponse{}, err
	}

	err = m.paymentService.CancelSubscription(req.SubscriptionID)

	type CreateSubscriptionRes struct {
		Created        bool   `json:"created"`
		SubscriptionId string `json:"subscription_id"`
	}
	res := CreateSubscriptionRes{Created: true}
	if err != nil {
		res.Created = false
	}
	return BuildRes(res, "Successfully created!", true)

}
