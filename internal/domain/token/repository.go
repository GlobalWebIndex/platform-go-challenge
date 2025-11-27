package token

import (
	"github.com/google/uuid"
	"time"
)

type RefreshRecord struct {
	UserID uuid.UUID
	Expiry time.Time
	Roles  []string
}
