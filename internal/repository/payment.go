package repository

import (
	"ownify_api/internal/dto"
	"ownify_api/internal/utils"
	"time"
)

type PaymentQuery interface {
	CreateCustomer(email, customerId string) error
	CreateSubscription(subscription dto.Subscription) error

	//email, customerId, priceId, subscriptionId string, endAt int64

	UpdateSubscription(customerId, priceId, subscriptionId string, endAt int64) error
	CancelSubscription(email string, customerId string) error
	VerifySubscriptionStatus(email string) bool
}

type paymentQuery struct{}

func (l *paymentQuery) CreateCustomer(email, customerId string) error {

	tableName := PaymentTableName
	sqlBuilder := utils.NewSqlBuilder()
	cols := []string{"email", "customer_id"}
	vals := []string{email, customerId}

	// Convert []string to []interface{}
	interfaceVals := make([]interface{}, len(vals))
	for i, v := range vals {
		interfaceVals[i] = v
	}

	sql, err := sqlBuilder.Insert(tableName, cols, interfaceVals)
	if err != nil {
		return err
	}
	_, err = DB.Exec(*sql)
	if err != nil {
		return err
	}
	return nil
}

func (l *paymentQuery) CreateSubscription(subscription dto.Subscription) error {

	tableName := PaymentTableName
	sqlBuilder := utils.NewSqlBuilder()
	cons, values := utils.ConvertToEntity(&subscription)
	sql, err := sqlBuilder.Insert(tableName, cons, values)
	if err != nil {
		return err
	}
	_, err = DB.Exec(*sql)
	if err != nil {
		return err
	}
	return nil
}

func (l *paymentQuery) UpdateSubscription(customerId string, priceId string, subscriptionId string, endAt int64) error {

	tableName := PaymentTableName
	sqlBuilder := utils.NewSqlBuilder()
	sql, err := sqlBuilder.Update(tableName, []utils.Tuple{{Key: "price_id", Val: priceId}, {Key: "subscription_id", Val: subscriptionId}, {Key: "end_at", Val: endAt}}, []utils.Tuple{{Key: "customer_id", Val: customerId}}, "OR")
	if err != nil {
		return err
	}
	_, err = DB.Exec(*sql)
	if err != nil {
		return err
	}
	return nil
}

func (l *paymentQuery) CancelSubscription(email string, customerId string) error {

	tableName := PaymentTableName
	sqlBuilder := utils.NewSqlBuilder()
	sql, err := sqlBuilder.Update(tableName, []utils.Tuple{{Key: "subscription_id", Val: ""}, {Key: "end_at", Val: 0}}, []utils.Tuple{{Key: "customer_id", Val: customerId}, {Key: "email", Val: email}}, "AND")
	if err != nil {
		return err
	}
	_, err = DB.Exec(*sql)
	if err != nil {
		return err
	}
	return nil
}

func (l *paymentQuery) VerifySubscriptionStatus(email string) bool {

	tableName := PaymentTableName
	sqlBuilder := utils.NewSqlBuilder()
	var customerId string
	var subscriptionId string
	var priceId string
	var endAt string
	sql, err := sqlBuilder.Select(tableName, []string{"customer_id", "subscription_id", "price_id", "end_at"}, []utils.Tuple{{Key: "email", Val: email}}, "=", "OR")
	if err != nil {
		return false
	}
	err = DB.QueryRow(*sql).Scan(
		&customerId,
		&subscriptionId,
		&priceId,
		&endAt,
	)
	if err != nil {
		return false
	}
	if customerId == "" || subscriptionId == "" || priceId == "" || endAt == "" {
		return false
	}
	//2023-06-13 09:15:50
	const layout = "2006-01-02 15:04:05"
	t, err := time.Parse(layout, endAt)
	if err != nil {
		return false
	}
	now := time.Now()
	return !t.After(now)
}
