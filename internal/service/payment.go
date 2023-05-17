package service

// import (
// 	"ownify_api/internal/repository"

// 	"github.com/stripe/stripe-go/v74"
// 	"github.com/stripe/stripe-go/v74/customer"
// )

// type PaymentService interface {
// 	CreatePlan(productName, currency string, priceUnitAmount int64, interval stripe.PriceRecurringInterval) (*stripe.Price, error)
// 	CreatePayment(amount int64, currency, paymentMethod string) (*stripe.PaymentIntent, error)
// }

// type paymentService struct {
// 	dbHandler    repository.DBHandler
// 	stripeAPIKey string
// }

// func NewPaymentService(dbHandler repository.DBHandler, stripeAPIKey string) PaymentService {
// 	stripe.Key = stripeAPIKey
// 	return &paymentService{
// 		dbHandler:    dbHandler,
// 		stripeAPIKey: stripeAPIKey,
// 	}
// }

// func (s *paymentService) CreatePlan(productName, currency string, priceUnitAmount int64, interval stripe.PriceRecurringInterval) (*stripe.Price, error) {
// 	p, err := product.New(&stripe.ProductParams{
// 		Name: stripe.String(productName),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	pr, err := price.New(&stripe.PriceParams{
// 		UnitAmount: stripe.Int64(priceUnitAmount),
// 		Currency:   stripe.String(currency),
// 		Recurring: &stripe.PriceRecurringParams{
// 			Interval: &interval,
// 		},
// 		Product: stripe.String(p.ID),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	// You can save product and price ID to your database here if needed

// 	return pr, nil
// }

// func (s *paymentService) CreatePayment(amount int64, currency, paymentMethod string) (*stripe.PaymentIntent, error) {
// 	params := &stripe.PaymentIntentParams{
// 		Amount:        stripe.Int64(amount),
// 		Currency:      stripe.String(currency),
// 		PaymentMethod: stripe.String(paymentMethod),
// 		Confirm:       stripe.Bool(true),
// 	}

// 	pi, err := paymentintent.New(params)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// You can save payment intent ID to your database here if needed

// 	return pi, nil
// }
