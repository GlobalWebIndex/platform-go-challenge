package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// User is an assets/favorites user
type User struct {
	Email    string
	Password string
}

// CustomClaims are the custom claims for the generated tokens
type CustomClaims struct {
	Email string `json:"email"`
}

// Claims are all the JWT claims for the generated tokens
type Claims struct {
	jwt.RegisteredClaims
	CustomClaims
}

// Tokenizer encodes/decodes JWT tokens.
// It uses symmetric encryption with HS256
type Tokenizer struct {
	secret           []byte               // the secret for the tokens
	keepAlive        time.Duration        // for token expiration
	now              func() time.Time     // a func, helps with testing
	registeredClaims jwt.RegisteredClaims // standard claims for all tokens
}

// AuthKey is the key type for Context
type AuthKey int
var (
	keyToken = AuthKey(0)
	keyUser = AuthKey(1)
)

// AuthService provides handlers for authenticating users
type AuthService struct {
	tokenizer *Tokenizer
}

var (
	// ErrTokenMissingCustomClaim is returned if a token does not include a custom claim
	ErrTokenMissingCustomClaim = errors.New("token is missing a custom claim")
)

// defaultTokenizer is a Tokenizer configured with the application defaults.
func defaultTokenizer() *Tokenizer {
	return &Tokenizer{
		secret:    configuredSecret,
		keepAlive: 8 * time.Hour,
		now:       func() time.Time { return time.Now() },
		registeredClaims: jwt.RegisteredClaims{
			Issuer:   "gwifav",
			Subject:  "favorites",
			Audience: []string{"users"},
		},
	}
}

// tokenForUser returns a JWT token for a user
func (t *Tokenizer) tokenForUser(u User) string {
	now := t.now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.registeredClaims.Issuer,
			Subject:   t.registeredClaims.Subject,
			Audience:  t.registeredClaims.Audience,
			ExpiresAt: jwt.NewNumericDate(now.Add(t.keepAlive)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		CustomClaims: CustomClaims{
			Email: u.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(t.secret)
	if err != nil {
		// this can happen if it fails to encode the claims as json
		// Because such cases should be caught on testing
		// it is better to stop here
		panic(fmt.Errorf("failed to create jwt token for user %s: %v", u.Email, err))
	}
	return tokenStr
}

// userFromToken extracts the user from the token.
// Not all fields of User are included (ex Password).
// The caller should check the user against the User DB
func (t *Tokenizer) userFromToken(rawToken string) (User, error) {
	parser := jwt.NewParser(
		jwt.WithExpirationRequired(),
		jwt.WithAudience(strings.Join(t.registeredClaims.Audience, " ")),
		jwt.WithIssuedAt(),
		jwt.WithIssuer(t.registeredClaims.Issuer),
		jwt.WithLeeway(1*time.Minute),
		jwt.WithSubject(t.registeredClaims.Subject),
		jwt.WithValidMethods([]string{"HS256"}),
	)

	token, err := parser.ParseWithClaims(rawToken, &Claims{}, func(_ *jwt.Token) (any, error) {
		return t.secret, nil
	})
	if err != nil {
		return User{}, fmt.Errorf("failed to extract user from JWT token: %w", err)
	}

	email := token.Claims.(*Claims).Email
	if email == "" {
		return User{}, fmt.Errorf("missing email claim: %w", ErrTokenMissingCustomClaim)
	}

	return User{Email: email}, nil
}

// loginHandler takes the user email and password, verifies them and returns a JWT token.
func (ah *AuthService) loginHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := strings.TrimSpace(r.PostFormValue("email"))
	userPasswd := strings.TrimSpace(r.PostFormValue("passwd"))
	user, err := newUser(userEmail + " " + userPasswd)
	if err != nil {
		http.Error(w, "bad input, please check email, password and try again", http.StatusBadRequest)
		return
	}

	user, present := findUser(userEmail, userPasswd)
	if !present {
		http.Error(w, "unauthorized, please check credentials", http.StatusUnauthorized)
		return
	}

	log.Printf("successful login for user %q", user.Email)
	token := ah.tokenizer.tokenForUser(user)
	w.Write([]byte(token))
}

// verifyTokenHandler is used with withAuthorization to verify the JWT token
func (ah *AuthService) verifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

// withAuthorization checks if the request contains a valid JWT token
func (ah *AuthService) withAuthorization(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawToken, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if rawToken == "" {
			http.Error(w, "unauthorized, please use a token", http.StatusUnauthorized)
			return
		}

		user, err := ah.tokenizer.userFromToken(rawToken)
		if err != nil {
			log.Println("JWT validation error:", err)
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Error(w, "token expired, please issue a new one", http.StatusUnauthorized)
			} else {
				http.Error(w, "unauthorized, please check the token", http.StatusUnauthorized)
			}
			return
		}

		// TODO: should we check user here?

		rctx := context.WithValue(r.Context(), keyUser, user)
		rctx = context.WithValue(rctx, keyToken, rawToken)
		h.ServeHTTP(w, r.WithContext(rctx))
	})
}
