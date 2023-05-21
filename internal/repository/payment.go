package repository

import (
	"ownify_api/internal/utils"
)

type PaymentQuery interface {
	CreateCustomer(email, customerId string) error
	UpdateSubscription(customerId string, priceId string, subscriptionId string, endAt int64) error
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
	sql, err := sqlBuilder.Select(tableName, []string{"customer_id", "subscription_id", "price_id"}, []utils.Tuple{{Key: "email", Val: email}}, "eq", "OR")
	if err != nil {
		return false
	}
	err = DB.QueryRow(*sql).Scan(
		&customerId,
		&subscriptionId,
		&priceId,
	)
	if err != nil {
		return false
	}
	if customerId == "" || subscriptionId == "" || priceId == "" {
		return false
	}
	return true
}
