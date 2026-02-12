package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ---------- helpers ----------

func unsignedToken(sub string, exp time.Time) string {
	claims := jwt.MapClaims{"sub": sub, "exp": exp.Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	s, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	return s
}

func unsignedTokenWithoutSubjectClaim(exp time.Time) string {
	claims := jwt.MapClaims{"exp": exp.Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	s, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	return s
}

func signedToken(sub, secret string, exp time.Time) string {
	claims := jwt.MapClaims{"sub": sub, "exp": exp.Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString([]byte(secret))
	return s
}

// dummyHandler writes 200 and the extracted userID.
var dummyHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	uid := UserIDFromContext(r.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(uid))
})

// ---------- tests ----------

func TestJWTMiddleware_UnsignedModeAllowed(t *testing.T) {
	mw := JWTMiddleware(AuthConfig{Secret: "", AllowUnsignedTokens: true})

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "valid unsigned token",
			authHeader: "Bearer " + unsignedToken("user42", time.Now().Add(time.Hour)),
			wantStatus: http.StatusOK,
			wantBody:   "user42",
		},
		{
			name:       "missing Authorization header",
			authHeader: "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "malformed header (no Bearer prefix)",
			authHeader: "Token abc",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "signed token rejected in unsigned mode",
			authHeader: "Bearer " + signedToken("user1", "secret", time.Now().Add(time.Hour)),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "garbage token",
			authHeader: "Bearer not.a.token",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "token missing sub claim",
			authHeader: "Bearer " + unsignedTokenWithoutSubjectClaim(time.Now().Add(time.Hour)),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "token with empty sub claim",
			authHeader: "Bearer " + unsignedToken("", time.Now().Add(time.Hour)),
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()
			mw(dummyHandler).ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", rr.Code, tt.wantStatus, rr.Body.String())
			}
			if tt.wantBody != "" && rr.Body.String() != tt.wantBody {
				t.Errorf("body = %q, want %q", rr.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestJWTMiddleware_UnsignedModeDisabled(t *testing.T) {
	// When AllowUnsignedTokens is false and no secret, all requests should be rejected
	mw := JWTMiddleware(AuthConfig{Secret: "", AllowUnsignedTokens: false})

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
	}{
		{
			name:       "unsigned token rejected when not allowed",
			authHeader: "Bearer " + unsignedToken("user42", time.Now().Add(time.Hour)),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "signed token also rejected (no secret configured)",
			authHeader: "Bearer " + signedToken("user1", "secret", time.Now().Add(time.Hour)),
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()
			mw(dummyHandler).ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", rr.Code, tt.wantStatus, rr.Body.String())
			}
		})
	}
}

func TestJWTMiddleware_SignedMode(t *testing.T) {
	const secret = "test-secret"
	mw := JWTMiddleware(AuthConfig{Secret: secret, AllowUnsignedTokens: false})

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "valid signed token",
			authHeader: "Bearer " + signedToken("user7", secret, time.Now().Add(time.Hour)),
			wantStatus: http.StatusOK,
			wantBody:   "user7",
		},
		{
			name:       "wrong secret",
			authHeader: "Bearer " + signedToken("user7", "wrong", time.Now().Add(time.Hour)),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "unsigned token rejected in signed mode",
			authHeader: "Bearer " + unsignedToken("user7", time.Now().Add(time.Hour)),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "expired token",
			authHeader: "Bearer " + signedToken("user7", secret, time.Now().Add(-time.Hour)),
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()
			mw(dummyHandler).ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", rr.Code, tt.wantStatus, rr.Body.String())
			}
			if tt.wantBody != "" && rr.Body.String() != tt.wantBody {
				t.Errorf("body = %q, want %q", rr.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestUserIDFromContext_EmptyWhenNoMiddleware(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	if uid := UserIDFromContext(req.Context()); uid != "" {
		t.Errorf("expected empty user ID, got %q", uid)
	}
}
