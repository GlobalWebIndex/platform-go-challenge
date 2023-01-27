package repository

import (
	//"fmt"

	"fmt"
	"ownify_api/internal/domain"
	"ownify_api/internal/dto"
	"ownify_api/internal/utils"

	//"ownify_api/internal/dto"

	sq "github.com/Masterminds/squirrel"
	//"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/status"
)

type BusinessQuery interface {
	CreateBusiness(
		business *dto.BriefBusiness,
	) error
	GetBusiness(email string) (*interface{}, error)
	DeleteBusiness(email string, userId string) error

	VerifyBusiness(userId string, email string) (*interface{}, error)
}

type businessQuery struct{}

func (u *businessQuery) CreateBusiness(
	business *dto.BriefBusiness,
) error {
	if !business.Valid() {
		return fmt.Errorf("[ERR] invalid Info: %v", business)
	}
	tableName := domain.BusinessTableName

	//cols1 := []string{"email", "pin", "business", "first_name", "last_name", "location", "phone_number"}
	//values1 := []interface{}{business.Email, business.Pin, business.Business, business.FirstName, business.LastName, business.Location, business.PhoneNumber}

	cols, values := utils.ConvertToEntity(business)
	sqlBuilder := utils.NewSqlBuilder()
	//fmt.Println(cols1)
	//fmt.Println(values1...)

	query, err := sqlBuilder.Insert(tableName, cols, values)
	//query1, _ := sqlBuilder.Insert(tableName, cols1, values1)
	fmt.Println(query)
	//fmt.Println(query1)

	if err != nil {
		return err
	}
	_, err = DB.Exec(*query)
	if err != nil {
		return err
	}
	return nil
}

func (b *businessQuery) GetBusiness(email string) (*interface{}, error) {
	var user interface{}
	err := pgQb().Select("*").Where(sq.Eq{"email": email}).From(domain.BusinessTableName).QueryRow().Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (b *businessQuery) DeleteBusiness(email string, userId string) error {
	sqlBuilder := utils.NewSqlBuilder()
	conditions := []utils.Tuple{{Key: "user_id", Val: userId}, {Key: "email", Val: email}}
	sql, err := sqlBuilder.Delete(domain.BusinessTableName, conditions, "AND")
	if err != nil {
		return err
	}
	_, err = DB.Exec(*sql)
	if err != nil {
		return err
	}
	return nil
}

func (b *businessQuery) VerifyBusiness(userId string, email string) (*interface{}, error) {
	var user interface{}
	sqlBuilder := utils.NewSqlBuilder()
	conditions := []utils.Tuple{{Key: "user_id", Val: userId}, {Key: "email", Val: email}}
	sql, err := sqlBuilder.Select(domain.BusinessTableName, []string{}, conditions, "AND")
	fmt.Println(*sql)
	if err != nil {
		return nil, err
	}

	user, err = DB.Exec(*sql)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
