package outbound

import "github.com/Kercyn/crud_template/internal/core/domain"

type UserRepository interface {
	GetByUserID(userID DataSourceID) (domain.User, error)
	Patch(request PatchData) error
}

type PatchData struct {
	ID     DataSourceID
	Fields map[string]interface{}
}
