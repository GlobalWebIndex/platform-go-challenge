package service

import (
	"fmt"
	"log"
	"ownify_api/internal/dto"
	"ownify_api/internal/repository"

	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/product"
	sub "github.com/stripe/stripe-go/v74/subscription"
)

type PaymentService interface {
	CreateProduct(name, price, description string) (string, error)
	CreateCheckoutSessionId(priceId string) (string, string, error)
	CreateCustomer(email string) (string, error)
	CreateSubscription(customerID string, priceId string) (*string, error)
	UpdateSubscription(oldPriceID, newPriceId string) error
	CancelSubscription(subscriptionID string) error

	VerifySubscriptionStatus(email string) bool
	GetActiveProducts() ([]*stripe.Product, error)
	CheckoutSession(sessionId string) error
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

func (s *paymentService) CreateProduct(name, price, description string) (string, error) {
	params := &stripe.ProductParams{
		Name:         stripe.String(name),
		Description:  stripe.String(description),
		Type:         stripe.String("basic service"), //or "good", according to your need
		DefaultPrice: &price,
	}
	p, err := product.New(params)
	if err != nil {
		return "", err
	}
	return p.ID, nil
}

func (s *paymentService) CreateCheckoutSessionId(priceId string) (string, string, error) {

	viper.AddConfigPath("../config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("cannot read from a config")
	}
	ownfiyUrl := viper.Get("ownify.client.url").(string)
	successUrl := ownfiyUrl + "subscription?session_id={CHECKOUT_SESSION_ID}"
	cancelUrl := ownfiyUrl + "subscription"
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceId),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(successUrl),
		CancelURL:  stripe.String(cancelUrl),
	}

	checkoutSession, err := session.New(params)
	if err != nil {
		return "", "", err
	}
	return checkoutSession.ID, checkoutSession.URL, nil
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

func (s *paymentService) GetActiveProducts() ([]*stripe.Product, error) {
	params := &stripe.ProductListParams{
		Active: stripe.Bool(true),
	}

	i := product.List(params)
	var products []*stripe.Product
	for i.Next() {
		p := i.Product()
		products = append(products, p)
	}

	if err := i.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *paymentService) CheckoutSession(sessionId string) error {
	// Retrieve the checkout session
	sess, err := session.Get(sessionId, nil)
	if err != nil {
		return err
	}

	// Get the customer object
	customerParams := &stripe.CustomerParams{}
	customerParams.AddExpand("subscriptions")
	cust, err := customer.Get(sess.Customer.ID, customerParams)
	if err != nil {
		return err
	}

	if s.dbHandler.NewPaymentQuery().VerifySubscriptionStatus(cust.Email) {
		return fmt.Errorf("[Err] Already registered this subscription")
	}

	subscriptionParams := &stripe.SubscriptionParams{}
	subscriptionParams.AddExpand("items.data.price")
	sub, err := sub.Get(sess.Subscription.ID, subscriptionParams)
	if err != nil {
		return err
	}

	customerID := cust.ID
	subscriptionID := sub.ID
	endedAt := sub.EndedAt
	priceID := sub.Items.Data[0].Price.ID

	subscription := dto.Subscription{
		Email:          cust.Email,
		CustomerId:     customerID,
		SubscriptionId: subscriptionID,
		EndAt:          endedAt,
		PriceId:        priceID,
	}

	err = s.dbHandler.NewPaymentQuery().CreateSubscription(subscription)
	return err
}
