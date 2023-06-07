package app

import (
	"context"
	"fmt"

	"ownify_api/internal/constants"
	desc "ownify_api/pkg"

	"github.com/stripe/stripe-go/v74"
)

func (m *MicroserviceServer) GetSubscriptionPlans(ctx context.Context, req *desc.GetSubscriptionPlansRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	_, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, err
	}

	plans, err := m.paymentService.GetActiveProducts()
	if err != nil {
		return nil, err
	}
	type CreateCheckoutSessionIdRes struct {
		Plans []*stripe.Product `json:"plans"`
	}

	res := CreateCheckoutSessionIdRes{
		Plans: plans,
	}
	return BuildRes(res, "Successfully created!", true)

}

func (m *MicroserviceServer) CreateCheckoutSessionId(ctx context.Context, req *desc.CreateCheckoutSessionIdRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	_, err := m.TokenInterceptor(ctx)
	if err != nil {
		return &desc.NetWorkResponse{}, err
	}

	sessionId, redirectUrl, err := m.paymentService.CreateCheckoutSessionId(req.PriceId)

	type CreateCheckoutSessionIdRes struct {
		Created     bool   `json:"created"`
		SessionId   string `json:"session_id"`
		RedirectUrl string `json:"redirect_url"`
	}
	res := CreateCheckoutSessionIdRes{Created: true, SessionId: sessionId, RedirectUrl: redirectUrl}
	if err != nil {
		res.Created = false
	}
	return BuildRes(res, "Successfully created!", true)
}

func (m *MicroserviceServer) CreateSubscription(ctx context.Context, req *desc.CreateSubscriptionRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	_, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUser, "raw message:%s", err)
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
		return nil, fmt.Errorf(constants.ErrInvalidUser, "raw message:%s", err)
	}
	err = m.paymentService.UpdateSubscription(req.OldPriceId, req.NewPriceId)

	type UpdateSubscriptionRes struct {
		Updated        bool   `json:"updated"`
		SubscriptionId string `json:"subscription_id"`
	}
	res := UpdateSubscriptionRes{Updated: true}
	if err != nil {
		res.Updated = false
	}
	return BuildRes(res, "Successfully created!", true)

}

func (m *MicroserviceServer) CancelSubscription(ctx context.Context, req *desc.CancelSubscriptionRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	_, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUser, "raw message:%s", err)
	}

	err = m.paymentService.CancelSubscription(req.SubscriptionID)

	type CancelSubscriptionRes struct {
		Canceled       bool   `json:"canceled"`
		SubscriptionId string `json:"subscription_id"`
	}
	res := CancelSubscriptionRes{Canceled: true}
	if err != nil {
		res.Canceled = false
	}
	return BuildRes(res, "Successfully created!", true)

}

func (m *MicroserviceServer) CheckSessionID(ctx context.Context, req *desc.CheckSessionIDRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	_, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUser, "raw message:%s", err)
	}

	err = m.paymentService.CheckoutSession(req.SessionId)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidStripSessionID, "raw message:%s", err)
	}

	return BuildRes(true, "Successfully session checked out!", true)
}
