package service

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/golang-jwt/jwt"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenManager interface {
	NewJWT(userID string) (string, error)
	Parse(accessToken string) (*int64, error)
	ParseFirebaseToken(accessToken string) (*int64, error)
	NewRefreshToken() (string, error)
	NewFirebaseToken(accessToken string, userId int64) (string, error)
	ValidateFirebase(accessToken string) (*string, error)
}

type tokenManager struct {
	signingKey string
}

func NewTokenManager(signedKey string) TokenManager {
	return &tokenManager{signingKey: signedKey}
}

func (t *tokenManager) NewJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		Subject:   userID,
	})
	return token.SignedString([]byte(t.signingKey))
}

func (t *tokenManager) NewFirebaseToken(accessToken string, userId int64) (string, error) {
	uid, err := t.ValidateFirebase(accessToken)
	if err != nil {
		return "", err
	}
	client, _ := getFirebaseAuthService()
	devClaims := make(map[string]interface{})
	devClaims["userId"] = userId
	token, err := client.CustomTokenWithClaims(context.Background(), *uid, devClaims)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (t *tokenManager) Parse(accessToken string) (*int64, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Errorf(codes.InvalidArgument, "unexpected signing method")
		}
		return []byte(t.signingKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("cannot get claims from token")
	}
	atoi, err := strconv.Atoi(claims["sub"].(string))
	if err != nil {
		return nil, fmt.Errorf("cannot convert str to int: %v", err)
	}
	id := int64(atoi)
	return &id, nil
}

func (t *tokenManager) ParseFirebaseToken(accessToken string) (*int64, error) {
	client, err := getFirebaseAuthService()
	if err != nil {
		return nil, err
	}
	token, err := client.VerifyIDToken(context.Background(), accessToken)
	if err != nil {
		return nil, err
	}
	userId, ok := token.Claims["userId"]
	if !ok {
		return nil, fmt.Errorf("cannot get claims from token")
	}
	atoi, err := strconv.Atoi(userId.(string))
	if err != nil {
		return nil, fmt.Errorf("cannot convert str to int: %v", err)
	}
	id := int64(atoi)
	return &id, nil
}

func (t *tokenManager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func (t *tokenManager) ValidateFirebase(accessToken string) (*string, error) {
	client, err := getFirebaseAuthService()
	if err != nil {
		return nil, fmt.Errorf("[ERR] can't initialize firebase app client : %s", err)
	}
	token, err := client.VerifyIDTokenAndCheckRevoked(context.Background(), accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid user: %s", err)
	}
	return &token.UID, nil
}

func getFirebaseAuthService() (*auth.Client, error) {
	opt := option.WithCredentialsFile("../config/ownify-wallet-service-account.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("[ERR] can't initialize firebase app: %s", err)
	}
	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("[ERR] can't initialize firebase app client : %s", err)
	}
	return client, nil
}
