package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/GlobalWebIndex/platform-go-challenge/internal/auth"
)

// stubAuthenticator is a test double implementing auth.Authenticator.
type stubAuthenticator struct {
	userID string
	err    error
}

func (s stubAuthenticator) Authenticate(r *http.Request) (string, error) {
	return s.userID, s.err
}

func TestAuthUserMiddleware_MissingAuth_Returns401(t *testing.T) {
	t.Parallel()

	stub := stubAuthenticator{
		userID: "",
		err:    auth.ErrUnauthenticated,
	}

	var nextCalled bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/user-1/favourites", nil)

	mw := NewAuthUserMiddleware(stub)
	mw(next).ServeHTTP(rr, req)

	if nextCalled {
		t.Fatalf("expected next handler not to be called when unauthenticated")
	}
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rr.Code)
	}
}

func TestAuthUserMiddleware_UserMismatch_Returns403(t *testing.T) {
	t.Parallel()

	stub := stubAuthenticator{
		userID: "user-1",
		err:    nil,
	}

	var nextCalled bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/other-user/favourites", nil)

	// Attach chi route context so URLParam("userId") works.
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userId", "other-user")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	mw := NewAuthUserMiddleware(stub)
	mw(next).ServeHTTP(rr, req)

	if nextCalled {
		t.Fatalf("expected next handler not to be called on user mismatch")
	}
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rr.Code)
	}
}

func TestAuthUserMiddleware_MatchingUser_AllowsRequestAndSetsContext(t *testing.T) {
	t.Parallel()

	stub := stubAuthenticator{
		userID: "user-1",
		err:    nil,
	}

	var (
		nextCalled bool
		gotUserID  string
		gotOK      bool
	)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		userID, ok := UserIDFromContext(r.Context())
		gotUserID, gotOK = userID, ok
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/user-1/favourites", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userId", "user-1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	mw := NewAuthUserMiddleware(stub)
	mw(next).ServeHTTP(rr, req)

	if !nextCalled {
		t.Fatalf("expected next handler to be called")
	}
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 from next handler, got %d", rr.Code)
	}
	if !gotOK || gotUserID != "user-1" {
		t.Fatalf("expected user ID %q in context, got %q (ok=%v)", "user-1", gotUserID, gotOK)
	}
}

func TestAuthUserMiddleware_UnexpectedError_Returns500(t *testing.T) {
	t.Parallel()

	stub := stubAuthenticator{
		userID: "",
		err:    errors.New("boom"),
	}

	var nextCalled bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/user-1/favourites", nil)

	mw := NewAuthUserMiddleware(stub)
	mw(next).ServeHTTP(rr, req)

	if nextCalled {
		t.Fatalf("expected next handler not to be called on unexpected auth error")
	}
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rr.Code)
	}
}
