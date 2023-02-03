package repository

import (
	//"fmt"

	"fmt"
	"ownify_api/internal/domain"
	"ownify_api/internal/dto"
	"ownify_api/internal/utils"
	//"ownify_api/internal/dto"
	//"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/status"
)

type BusinessQuery interface {
	CreateBusiness(
		business *dto.BriefBusiness,
	) error
	GetBusiness(email string) (*dto.BriefBusiness, error)
	GetBusinessByWalletAddress(pubKey string) (*dto.BriefBusiness, error)
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

func (b *businessQuery) GetBusiness(email string) (*dto.BriefBusiness, error) {
	var user dto.BriefBusiness

	sqlBuilder := utils.NewSqlBuilder()
	sql, err := sqlBuilder.Select(domain.BusinessTableName, []string{
		"business",
		"first_name",
		"last_name",
		"location",
		"phone_number",
		"user_id",
	}, []utils.Tuple{{Key: "email", Val: email}}, "=", "And")
	if err != nil {
		return nil, err
	}
	err = DB.QueryRow(*sql).Scan(
		&user.Business,
		&user.FirstName,
		&user.LastName,
		&user.Location,
		&user.PhoneNumber,
		&user.UserId,
	)
	user.Email = email
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (b *businessQuery) GetBusinessByWalletAddress(pubKey string) (*dto.BriefBusiness, error) {
	var user dto.BriefBusiness
	sql := fmt.Sprintf("SELECT b.email b.business b.first_name b.last_name b.location FROM %s w LEFT JOIN %s b ON w.email = b.email WHERE w.pub_addr = \"%s\"", domain.WalletTableName, domain.BusinessTableName, pubKey)
	err := DB.QueryRow(sql).Scan(
		&user.Email,
		&user.Business,
		&user.FirstName,
		&user.LastName,
		&user.Location,
	)
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
	sql, err := sqlBuilder.Select(domain.BusinessTableName, []string{}, conditions, "=", "AND")
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
