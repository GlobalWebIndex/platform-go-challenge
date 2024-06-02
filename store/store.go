package store

import (
	"context"
	"fmt"
	"os"
	"time"
	"user-favorites-api/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var JwtKey = []byte(os.Getenv("SECRET_KEY"))

type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

type Store struct {
    DB *pgxpool.Pool
}

func NewStore(connectionString string) (*Store, error) {
    config, err := pgxpool.ParseConfig(connectionString)
    if err != nil {
        return nil, err
    }
    db, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        return nil, err
    }

    return &Store{DB: db}, nil
}


func (s *Store) RegisterUser(user *models.User) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user.Password = string(hashedPassword)
    _, err = s.DB.Exec(context.Background(), "INSERT INTO users (username, email, password) VALUES ($1, $2, $3)",
        user.Username, user.Email, user.Password)
    return err
}

func (s *Store) AuthenticateUser(credentials *models.Credentials) (string, error) {
    var user models.User
    err := s.DB.QueryRow(context.Background(), "SELECT id, password FROM users WHERE username=$1", credentials.Username).Scan(&user.ID, &user.Password)
    if err != nil {
        return "", fmt.Errorf("user not found")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
        return "", fmt.Errorf("invalid credentials")
    }

    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Username: credentials.Username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(), 
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(JwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func (s *Store) GetUserFavorites(userID string) (*models.UserFavorites, error) {
    userFavorites := &models.UserFavorites{UserID: userID}
    rows, err := s.DB.Query(context.Background(), "SELECT asset_id, type, description, data FROM assets WHERE user_id=$1", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var asset models.Asset
        if err := rows.Scan(&asset.ID, &asset.Type, &asset.Description, &asset.Data); err != nil {
            return nil, err
        }
        userFavorites.Assets = append(userFavorites.Assets, asset)
    }

    return userFavorites, nil
}

func (s *Store) AddFavorite(userID string, asset models.Asset) error {
    _, err := s.DB.Exec(context.Background(), "INSERT INTO assets (user_id, asset_id, type, description, data) VALUES ($1, $2, $3, $4, $5)",
        userID, asset.ID, asset.Type, asset.Description, asset.Data)
    return err
}

func (s *Store) RemoveFavorite(userID, assetID string) error {
    _, err := s.DB.Exec(context.Background(), "DELETE FROM assets WHERE user_id=$1 AND asset_id=$2", userID, assetID)
    return err
}

func (s *Store) EditFavorite(userID, assetID string, asset models.Asset) error {
    _, err := s.DB.Exec(context.Background(), "UPDATE assets SET description=$1, data=$2 WHERE user_id=$3 AND asset_id=$4",
        asset.Description, asset.Data, userID, assetID)
    return err
}
