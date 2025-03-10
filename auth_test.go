package main

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func TestAuthTokenizer(t *testing.T) {
	t.Run("encodeDecode", func(t *testing.T) {
		userE := User{"userA@dev", "passwd"}
		token := mustGetUserToken(t, userE, defaultTokenizer())
		userD, err := defaultTokenizer().userFromToken(token)
		if err != nil {
			t.Errorf("failed to decode token: %v", err)
			t.FailNow()
		}
		if userD.Email != userE.Email {
			t.Errorf("email mismatch: got %s expected %s", userD.Email, userE.Email)
		}
	})

	t.Run("decodeFromOtherProducer", func(t *testing.T) {
		// a token generated from http://jwtbuilder.jamiekurtz.com/ signed with "secret"
		// {
		//   "iss": "gwifav",
		//   "iat": 1741273994,
		//   "exp": 1772810003,
		//   "aud": "users",
		//   "sub": "favorites",
		//   "Email": "jrocket@example.com"
		// }
		rawToken := `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJnd2lmYXYiLCJpYXQiOjE3NDEyNzM5OTQsImV4cCI6MTc3MjgxMDAwMywiYXVkIjoidXNlcnMiLCJzdWIiOiJmYXZvcml0ZXMiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20ifQ.uxeBJgDvNQGSki6YWPRMkG8e-1JpnhpYgqOGI_96t1A`
		tokenizer := &Tokenizer{
			secret: []byte("secret"),
			registeredClaims: jwt.RegisteredClaims{
				Issuer:   "gwifav",
				Subject:  "favorites",
				Audience: []string{"users"},
			},
		}

		user, err := tokenizer.userFromToken(rawToken)
		if err != nil {
			t.Errorf("failed to decode token: %v", err)
			t.FailNow()
		}
		if user.Email != "jrocket@example.com" {
			t.Errorf("email mismatch: got %v expected %v", user.Email, "jrocket@example.com")
			t.FailNow()
		}
	})

	t.Run("decodeFromDifferentAlgorithm", func(t *testing.T) {
		now := time.Now()
		claims := Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "gwifav",
				Subject:   "favorites",
				Audience:  []string{"users"},
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
			},
			CustomClaims: CustomClaims{
				Email: "user@dev",
			},
		}
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(configuredSecret)
		if err != nil {
			t.Errorf("failed to generate test token: %v", err)
			t.FailNow()
		}

		_, err = defaultTokenizer().userFromToken(token)
		if err == nil {
			t.Errorf("parser accepts tokens with a different signing algorithm")
		}
	})

	t.Run("decodeFromDifferentSecret", func(t *testing.T) {
		now := time.Now()
		claims := Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "gwifav",
				Subject:   "favorites",
				Audience:  []string{"users"},
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
			},
			CustomClaims: CustomClaims{
				Email: "user@dev",
			},
		}
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("other"))
		if err != nil {
			t.Errorf("failed to generate test token: %v", err)
			t.FailNow()
		}

		_, err = defaultTokenizer().userFromToken(token)
		if err == nil {
			t.Errorf("parser accepts tokens signed with a different secret")
		}
	})

	t.Run("decodeExpiredToken", func(t *testing.T) {
		now := time.Now()
		claims := Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "gwifav",
				Subject:   "favorites",
				Audience:  []string{"users"},
				ExpiresAt: jwt.NewNumericDate(now.Add(-1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
			},
			CustomClaims: CustomClaims{
				Email: "user@dev",
			},
		}
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(configuredSecret)
		if err != nil {
			t.Errorf("failed to generate test token: %v", err)
			t.FailNow()
		}

		_, err = defaultTokenizer().userFromToken(token)
		if err == nil {
			t.Errorf("parser accepts expired token")
		}
	})
}

func TestAuthorization(t *testing.T) {
	silenceLogger(t)
	prevUsers := users
	prevSecret := configuredSecret
	t.Cleanup(func() {
		users = prevUsers
		configuredSecret = prevSecret
	})
	userA := User{"userA@dev", "passwd"}
	userB := User{"userB@dev", "passwd"}
	users = []User{userA, userB}
	configuredSecret = []byte("secretA")

	ah := &AuthService{defaultTokenizer()}
	verifyHandler := ah.withAuthorization(http.HandlerFunc(ah.verifyTokenHandler))

	t.Run("withoutToken", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/verify", nil)
		w := httptest.NewRecorder()
		verifyHandler.ServeHTTP(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("auth handler accepted empty token")
		}
	})

	t.Run("withIncorrectToken", func(t *testing.T) {
		now := time.Now()

		// expired token
		token := mustGetUserToken(t, userA, &Tokenizer{
			secret:    configuredSecret,
			keepAlive: -1 * time.Hour,
			now:       func() time.Time { return now },
			registeredClaims: jwt.RegisteredClaims{
				Issuer:   "gwifav",
				Subject:  "favorites",
				Audience: []string{"users"},
			},
		})
		req := httptest.NewRequest("GET", "/verify", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		verifyHandler.ServeHTTP(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("auth handler accepted expired token")
		}
	})

	t.Run("withCorrectToken", func(t *testing.T) {
		token := mustGetUserToken(t, userA, &Tokenizer{
			secret:    configuredSecret,
			keepAlive: 8 * time.Hour,
			now:       func() time.Time { return time.Now() },
			registeredClaims: jwt.RegisteredClaims{
				Issuer:   "gwifav",
				Subject:  "favorites",
				Audience: []string{"users"},
			},
		})
		req := httptest.NewRequest("GET", "/verify", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		verifyHandler.ServeHTTP(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("auth handler rejected correct token")
		}
	})
}

func TestAuthLogin(t *testing.T) {
	silenceLogger(t)
	prevUsers := users
	prevSecret := configuredSecret
	t.Cleanup(func() {
		users = prevUsers
		configuredSecret = prevSecret
	})
	userA := User{"userA@dev", "passwd"}
	userB := User{"userB@dev", "passwd"}
	users = []User{userA, userB}
	configuredSecret = []byte("secretA")

	ah := &AuthService{defaultTokenizer()}
	loginHandler := http.HandlerFunc(ah.loginHandler)

	t.Run("withNoCredentials", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/login", nil)
		w := httptest.NewRecorder()
		loginHandler(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("login handler accepted empty credentials")
		}
	})

	t.Run("withUnknownUser", func(t *testing.T) {
		v := url.Values{}
		v.Set("email", "unknown@dev")
		v.Set("passwd", "passwd")
		req := httptest.NewRequest("POST", "/login", nil)
		req.PostForm = v
		w := httptest.NewRecorder()
		loginHandler(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("login handler accepted unknown user")
		}
	})

	t.Run("withWrongCredentials", func(t *testing.T) {
		v := url.Values{}
		v.Set("email", "userA@dev")
		v.Set("passwd", "wrong")
		req := httptest.NewRequest("POST", "/login", nil)
		req.PostForm = v
		w := httptest.NewRecorder()
		loginHandler(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("login handler accepted wrong credentials")
		}
	})

	t.Run("withKnownUser", func(t *testing.T) {
		v := url.Values{}
		v.Set("email", "userA@dev")
		v.Set("passwd", "passwd")
		req := httptest.NewRequest("POST", "/login", nil)
		req.PostForm = v
		w := httptest.NewRecorder()
		loginHandler(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("login handler failed to generate a token")
		}

		// parse it and validate claims
		rawToken, _ := io.ReadAll(resp.Body)
		user, err := defaultTokenizer().userFromToken(string(rawToken))
		if err != nil {
			t.Errorf("login handler returned an invalid token: %v", err)
		}
		if user.Email != "userA@dev" {
			t.Errorf("login handler email mismatch: got %s expected %s", user.Email, "gwifav@dev")
		}
	})
}

func TestAuthTokenClaims(t *testing.T) {
	now := time.Now().UTC()
	otherSecret := []byte("secret")
	oneHour := time.Hour

	var tests = []struct {
		claims         Claims
		secretProducer []byte
		err            error
	}{
		{
			addCustomClaims("user@dev", addExpirationClaims(now, oneHour, claims(jwt.RegisteredClaims{
				Issuer:   "gwifav",
				Subject:  "favorites",
				Audience: []string{"users"},
			}))), configuredSecret, nil,
		},
		{
			addExpirationClaims(now, oneHour, claims(jwt.RegisteredClaims{
				Issuer:   "gwifav",
				Subject:  "favorites",
				Audience: []string{"users"},
			})), configuredSecret, ErrTokenMissingCustomClaim,
		},
		{
			addExpirationClaims(now, oneHour, claims(jwt.RegisteredClaims{
				Issuer:   "gwifav",
				Subject:  "favorites",
				Audience: []string{"users"},
			})), otherSecret, jwt.ErrTokenSignatureInvalid,
		},
		{
			addExpirationClaims(now, oneHour, claims(jwt.RegisteredClaims{
				Issuer:   "unknown",
				Subject:  "favorites",
				Audience: []string{"users"},
			})), configuredSecret, jwt.ErrTokenInvalidIssuer,
		},
		{
			addExpirationClaims(now, oneHour, claims(jwt.RegisteredClaims{
				Issuer:   "gwifav",
				Subject:  "unknown",
				Audience: []string{"users"},
			})), configuredSecret, jwt.ErrTokenInvalidSubject,
		},
		{
			addExpirationClaims(now, oneHour, claims(jwt.RegisteredClaims{
				Issuer:   "gwifav",
				Subject:  "favorites",
				Audience: []string{"unknown"},
			})), configuredSecret, jwt.ErrTokenInvalidAudience,
		},
		{
			claims(jwt.RegisteredClaims{
				Issuer:   "gwifav",
				Subject:  "favorites",
				Audience: []string{"users"},
			}), configuredSecret, jwt.ErrTokenRequiredClaimMissing,
		},
		{
			claims(jwt.RegisteredClaims{
				Issuer:   "gwifav",
				Subject:  "favorites",
				Audience: []string{"users"},
				IssuedAt: jwt.NewNumericDate(now),
			}), configuredSecret, jwt.ErrTokenRequiredClaimMissing,
		},
		{
			claims(jwt.RegisteredClaims{
				Issuer:    "gwifav",
				Subject:   "favorites",
				Audience:  []string{"users"},
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				ExpiresAt: jwt.NewNumericDate(now.Add(-1 * oneHour)),
			}), configuredSecret, jwt.ErrTokenExpired,
		},
	}

	tokenizer := defaultTokenizer()
	for i, tst := range tests {
		caseNum := i + 1
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, tst.claims).SignedString(tst.secretProducer)
		if err != nil {
			t.Errorf("failed to generate test token for case %d: %v", caseNum, err)
			t.FailNow()
		}
		_, err = tokenizer.userFromToken(token)
		if !errors.Is(err, tst.err) {
			t.Errorf("failed case %d: got: %v expected: %v", caseNum, err, tst.err)
		}
	}
}

func claims(base jwt.RegisteredClaims) Claims {
	return Claims{RegisteredClaims: base}
}

func addExpirationClaims(now time.Time, exp time.Duration, base Claims) Claims {
	return Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    base.Issuer,
			Subject:   base.Subject,
			Audience:  base.Audience,
			ExpiresAt: jwt.NewNumericDate(now.Add(exp)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		CustomClaims: CustomClaims{
			Email: base.Email,
		},
	}
}

func addCustomClaims(email string, base Claims) Claims {
	return Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    base.Issuer,
			Subject:   base.Subject,
			Audience:  base.Audience,
			ExpiresAt: base.ExpiresAt,
			IssuedAt:  base.IssuedAt,
			NotBefore: base.NotBefore,
		},
		CustomClaims: CustomClaims{
			Email: email,
		},
	}
}

func mustGetUserToken(t *testing.T, user User, tokenizer *Tokenizer) string {
	return tokenizer.tokenForUser(user)
}
