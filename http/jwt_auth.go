package http

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type Signer interface {
	Sign(claims *jwt.RegisteredClaims) (string, error)
}

type Verifier interface {
	Verify(tokenString string) (subject string, err error)
}

type hmacAuth struct {
	secret []byte
}

func (h hmacAuth) Sign(claims *jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.secret)
}

func (h hmacAuth) Verify(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return h.secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return "", err
	}
	if claims.Subject == "" {
		return "", jwt.ErrTokenInvalidClaims
	}
	return claims.Subject, nil
}

func NewHMACAuth(secret []byte) (Signer, Verifier) {
	a := hmacAuth{secret: secret}
	return a, a
}

type rsaSigner struct {
	key *rsa.PrivateKey
}

func (r rsaSigner) Sign(claims *jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(r.key)
}

type rsaVerifier struct {
	key *rsa.PublicKey
}

func (r rsaVerifier) Verify(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return r.key, nil
	}, jwt.WithValidMethods([]string{"RS256"}))
	if err != nil {
		return "", err
	}
	if claims.Subject == "" {
		return "", jwt.ErrTokenInvalidClaims
	}
	return claims.Subject, nil
}

func NewRSAAuth(privateKeyPEM, publicKeyPEM []byte) (Signer, Verifier, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, nil, fmt.Errorf("parse private key: %w", err)
	}
	var publicKey *rsa.PublicKey
	if len(publicKeyPEM) > 0 {
		publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)
		if err != nil {
			return nil, nil, fmt.Errorf("parse public key: %w", err)
		}
	} else {
		publicKey = &privateKey.PublicKey
	}
	return rsaSigner{key: privateKey}, rsaVerifier{key: publicKey}, nil
}

// LoadJWTAuth uses JWT_PRIVATE_KEY_PATH (+ optional JWT_PUBLIC_KEY_PATH) for RSA, else JWT_SECRET for HMAC. Falls back to dev secret if unset.
func LoadJWTAuth() (Signer, Verifier, error) {
	privatePath := os.Getenv("JWT_PRIVATE_KEY_PATH")
	if privatePath != "" {
		privatePEM, err := os.ReadFile(privatePath)
		if err != nil {
			return nil, nil, fmt.Errorf("read private key: %w", err)
		}
		publicPath := os.Getenv("JWT_PUBLIC_KEY_PATH")
		var publicPEM []byte
		if publicPath != "" {
			publicPEM, err = os.ReadFile(publicPath)
			if err != nil {
				return nil, nil, fmt.Errorf("read public key: %w", err)
			}
		}
		return NewRSAAuth(privatePEM, publicPEM)
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-in-production"
	}
	signer, verifier := NewHMACAuth([]byte(secret))
	return signer, verifier, nil
}
