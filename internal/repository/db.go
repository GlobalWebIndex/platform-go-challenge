package repository

import (
	"fmt"
	"gwi_api/internal/domain"
	"gwi_api/internal/dto"

	_ "github.com/go-sql-driver/mysql"
)

type MemoryDB struct {
	nextUserID  uint64
	nextAssetID uint64
	users       map[uint64]domain.User
	userEmails  map[string]uint64
	assets      map[uint64]domain.Asset
}

func (d *MemoryDB) IncreaseUserID() {
	d.nextUserID += 1
}

func (d *MemoryDB) RegisterUser(user dto.UserDto) (*uint64, error) {
	_, found := DB.userEmails[user.Email]
	if found {
		return nil, fmt.Errorf("already exist user")
	}
	id := d.nextUserID
	DB.userEmails[user.Email] = id
	DB.users[id] = domain.User{
		ID:        id,
		Password:  user.Password,
		Email:     user.Email,
		Favorites: make([]domain.Asset, 0),
	}
	d.IncreaseUserID()
	return &id, nil
}

func (d *MemoryDB) DeleteUser(userId uint64) error {
	user, found := DB.users[userId]
	if !found {
		return fmt.Errorf("[Err] user don't exist")
	}
	delete(DB.userEmails, user.Email)
	delete(DB.users, userId)
	return nil
}

func (d *MemoryDB) IncreaseAssetID() {
	d.nextAssetID += 1
}

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		nextUserID:  0,
		nextAssetID: 0,
		users:       make(map[uint64]domain.User),
		userEmails:  make(map[string]uint64),
		assets:      make(map[uint64]domain.Asset),
	}
}

type DBHandler interface {
	NewUserQuery() UserQuery
}

type dbHandler struct {
	db *MemoryDB
}

var DB *MemoryDB

func NewDBHandler(db *MemoryDB) DBHandler {
	DB = db
	return &dbHandler{db}
}

func NewDB() (*MemoryDB, error) {
	DB := NewMemoryDB()
	return DB, nil
}

func (d *dbHandler) NewUserQuery() UserQuery {
	return &userQuery{}
}
