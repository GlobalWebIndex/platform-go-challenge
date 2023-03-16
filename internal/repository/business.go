package repository

import (
	//"fmt"

	"fmt"
	"ownify_api/internal/dto"
	"ownify_api/internal/utils"
)

type BusinessQuery interface {
	CreateBusiness(
		business *dto.BriefBusiness,
	) error
	UpdateBusiness(
		business *dto.BriefBusiness,
	) error
	GetBusiness(email string) (*dto.BriefBusiness, error)
	GetBusinessByWalletAddress(pubKey string) (*dto.BriefBusiness, error)
	GetBusinessByUserId(userId string) (*dto.BriefBusiness, error)
	DeleteBusiness(email string, userId string) error

	VerifyBusiness(userId string, email string) error
	VerifyBusinessByUserId(userId string) error
}

type businessQuery struct{}

func (u *businessQuery) CreateBusiness(
	business *dto.BriefBusiness,
) error {
	if !business.Valid() {
		return fmt.Errorf("[ERR] invalid Info: %v", business)
	}
	tableName := BusinessTableName
	cols, values := utils.ConvertToEntity(business)
	if !utils.Contains(cols, "phone_number") {
		cols = append(cols, "phone_number")
		values = append(values, "")
	}

	sqlBuilder := utils.NewSqlBuilder()

	query, err := sqlBuilder.Insert(tableName, cols, values)
	if err != nil {
		return err
	}
	_, err = DB.Exec(*query)
	if err != nil {
		return err
	}
	return nil
}

func (u *businessQuery) UpdateBusiness(
	business *dto.BriefBusiness,
) error {
	if !business.Valid() {
		return fmt.Errorf("[ERR] invalid Info: %v", business)
	}
	tableName := BusinessTableName
	cols, values := utils.ConvertToEntity(business)
	if !utils.Contains(cols, "phone_number") {
		cols = append(cols, "phone_number")
		values = append(values, "")
	}

	sqlBuilder := utils.NewSqlBuilder()

	deleteSql, _ := sqlBuilder.Delete(tableName, []utils.Tuple{{Key: "user_id", Val: business.UserId}}, "=")
	_, err := DB.Exec(*deleteSql)
	if err != nil {
		return err
	}

	query, err := sqlBuilder.Insert(tableName, cols, values)
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
	sql, err := sqlBuilder.Select(BusinessTableName, []string{
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
	emailSql := fmt.Sprintf("SELECT b.email,b.business,b.first_name,b.last_name,b.location FROM %s w LEFT JOIN %s b ON w.email = b.email OR w.user_id = b.user_id WHERE w.pub_addr = \"%s\"", WalletTableName, BusinessTableName, pubKey)
	err := DB.QueryRow(emailSql).Scan(
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

func getBusiness(pubKey string) (*dto.BriefBusiness, error) {
	var user dto.BriefBusiness
	emailSql := fmt.Sprintf("SELECT b.email,b.business,b.first_name,b.last_name,b.location FROM %s w LEFT JOIN %s b ON w.email = b.email OR w.user_id = b.user_id WHERE w.pub_addr = \"%s\"", WalletTableName, BusinessTableName, pubKey)
	err := DB.QueryRow(emailSql).Scan(
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

// func getBusinessWithUserId(pubKey string) (*dto.BriefBusiness, error) {
// 	var user dto.BriefBusiness
// 	emailSql := fmt.Sprintf("SELECT b.email,b.business,b.first_name,b.last_name,b.location FROM %s w LEFT JOIN %s b ON w.user_id = b.user_id WHERE w.pub_addr = \"%s\"", domain.WalletTableName, domain.BusinessTableName, pubKey)
// 	err := DB.QueryRow(emailSql).Scan(
// 		&user.Email,
// 		&user.Business,
// 		&user.FirstName,
// 		&user.LastName,
// 		&user.Location,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

func (b *businessQuery) GetBusinessByUserId(userId string) (*dto.BriefBusiness, error) {
	var user dto.BriefBusiness

	sqlBuilder := utils.NewSqlBuilder()
	sql, err := sqlBuilder.Select(BusinessTableName, []string{
		"email",
		"business",
		"first_name",
		"last_name",
		"location",
		"phone_number",
		"user_id",
	}, []utils.Tuple{{Key: "user_id", Val: userId}}, "=", "And")
	if err != nil {
		return nil, err
	}
	err = DB.QueryRow(*sql).Scan(
		&user.Email,
		&user.Business,
		&user.FirstName,
		&user.LastName,
		&user.Location,
		&user.PhoneNumber,
		&user.UserId,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (b *businessQuery) DeleteBusiness(email string, userId string) error {
	sqlBuilder := utils.NewSqlBuilder()
	conditions := []utils.Tuple{{Key: "user_id", Val: userId}, {Key: "email", Val: email}}
	sql, err := sqlBuilder.Delete(BusinessTableName, conditions, "AND")
	if err != nil {
		return err
	}
	_, err = DB.Exec(*sql)
	if err != nil {
		return err
	}
	return nil
}

func (b *businessQuery) VerifyBusiness(userId string, email string) error {
	isApproved := false
	sqlBuilder := utils.NewSqlBuilder()
	conditions := []utils.Tuple{{Key: "user_id", Val: userId}, {Key: "email", Val: email}}
	sql, err := sqlBuilder.Select(BusinessTableName, []string{"is_approved"}, conditions, "=", "AND")
	fmt.Println(*sql)
	if err != nil {
		return err
	}

	err = DB.QueryRow(*sql).Scan(&isApproved)
	if err != nil {
		return err
	}
	if isApproved {
		return fmt.Errorf("[ERR] already approved")
	}
	return nil
}

func (b *businessQuery) VerifyBusinessByUserId(userId string) error {
	isApproved := false
	sqlBuilder := utils.NewSqlBuilder()
	conditions := []utils.Tuple{{Key: "user_id", Val: userId}}
	sql, err := sqlBuilder.Select(BusinessTableName, []string{"is_approved"}, conditions, "=", "AND")
	fmt.Println(*sql)
	if err != nil {
		return err
	}

	err = DB.QueryRow(*sql).Scan(&isApproved)
	if err != nil {
		return err
	}
	if isApproved {
		return fmt.Errorf("[ERR] already approved")
	}
	return nil
}
