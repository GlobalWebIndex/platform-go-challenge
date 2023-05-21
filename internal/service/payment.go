package service

import (
	"ownify_api/internal/repository"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	sub "github.com/stripe/stripe-go/v74/subscription"
)

type PaymentService interface {
	CreateCustomer(email string) (string, error)
	CreateSubscription(customerID string, priceId string) (*string, error)
	UpdateSubscription(oldPriceID, newPriceId string) error
	CancelSubscription(subscriptionID string) error

	VerifySubscriptionStatus(email string) bool
}

type paymentService struct {
	dbHandler    repository.DBHandler
	stripeAPIKey string
}

func NewPaymentService(dbHandler repository.DBHandler, stripeAPIKey string) PaymentService {
	stripe.Key = stripeAPIKey
	return &paymentService{
		dbHandler:    dbHandler,
		stripeAPIKey: stripeAPIKey,
	}
}

func (s *paymentService) CreateCustomer(email string) (string, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}
	c, err := customer.New(params)
	if err != nil {
		return "", err
	}
	err = s.dbHandler.NewPaymentQuery().CreateCustomer(email, c.ID)
	if err != nil {
		return "", err
	}
	return c.ID, err
}

func (s *paymentService) CreateSubscription(customerID string, priceId string) (*string, error) {
	// Create subscription
	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceId),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
	}

	subscriptionParams.AddExpand("latest_invoice.payment_intent")
	newSub, err := sub.New(subscriptionParams)
	if err != nil {
		return &newSub.ID, err
	}
	err = s.dbHandler.NewPaymentQuery().UpdateSubscription(customerID, priceId, newSub.ID, newSub.EndedAt)
	return &newSub.ID, err
}

func (s *paymentService) UpdateSubscription(oldPriceID, newPriceId string) error {
	params := &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{{
			ID:    stripe.String(oldPriceID),
			Price: stripe.String(newPriceId),
		}},
	}

	updatedSubscription, err := sub.Update(oldPriceID, params)

	if err != nil {
		return err
	}
	err = s.dbHandler.NewPaymentQuery().UpdateSubscription(
		updatedSubscription.Customer.ID, newPriceId, updatedSubscription.ID, updatedSubscription.EndedAt)

	return err
}

func (s *paymentService) CancelSubscription(subscriptionID string) error {
	res, err := sub.Cancel(subscriptionID, nil)

	if err != nil {
		return err
	}

	err = s.dbHandler.NewPaymentQuery().CancelSubscription(
		res.Customer.Email, res.Customer.ID)
	return err
}

func (s *paymentService) VerifySubscriptionStatus(email string) bool {
	return s.dbHandler.NewPaymentQuery().VerifySubscriptionStatus(email)
}
