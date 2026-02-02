package httpapi_test

import (
	"bytes"
	"context"

	// "crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	assetrepo "github.com/manos/favourites/assets/adapters/outbound/repo_inmemory"
	assetapplication "github.com/manos/favourites/assets/application"
	"github.com/manos/favourites/config"
	"github.com/manos/favourites/favourites/adapters/outbound/assets_client"
	favrepo "github.com/manos/favourites/favourites/adapters/outbound/repo_inmemory"
	favouritesapplication "github.com/manos/favourites/favourites/application"
	httpapi "github.com/manos/favourites/http"
	"github.com/manos/favourites/http/auth"

	"github.com/golang-jwt/jwt/v5"
)

func TestFavouritesEndpointsHappyPath(t *testing.T) {
	ts, token := newTestServer(t)
	defer ts.Close()

	client := ts.Client()

	// initial list
	page := favouritesapplication.FavouritePageDTO{}
	doJSON(t, client, token, http.MethodGet, ts.URL+"/v1/favourites", nil, http.StatusOK, &page)
	if page.Total != 3 {
		t.Fatalf("initial favourites total = %d, want 3", page.Total)
	}

	// add favourite
	addBody := map[string]any{
		"type":        "INSIGHT",
		"assetId":     "ins-003",
		"description": "new fav",
	}
	var favID string
	doJSON(t, client, token, http.MethodPost, ts.URL+"/v1/favourites", addBody, http.StatusCreated, &favID)
	if favID == "" {
		t.Fatalf("expected favourite id, got empty")
	}

	// update description
	updateBody := map[string]any{"description": "updated desc"}
	doJSON(t, client, token, http.MethodPatch, ts.URL+"/v1/favourites/"+favID, updateBody, http.StatusNoContent, nil)

	// delete favourite
	doJSON(t, client, token, http.MethodDelete, ts.URL+"/v1/favourites/"+favID, nil, http.StatusNoContent, nil)

	// final list should be back to 3
	page = favouritesapplication.FavouritePageDTO{}
	doJSON(t, client, token, http.MethodGet, ts.URL+"/v1/favourites", nil, http.StatusOK, &page)
	if page.Total != 3 {
		t.Fatalf("final favourites total = %d, want 3", page.Total)
	}
}

func TestAssetsEndpointHappyPath(t *testing.T) {
	ts, token := newTestServer(t)
	defer ts.Close()

	client := ts.Client()

	var asset map[string]any
	doJSON(t, client, token, http.MethodGet, ts.URL+"/v1/assets/insights/ins-001", nil, http.StatusOK, &asset)
	if asset["id"] != "ins-001" {
		t.Fatalf("asset id = %v, want ins-001", asset["id"])
	}
}

func newTestServer(t *testing.T) (*httptest.Server, string) {
	t.Helper()

	cfg, err := config.LoadWithEnv("test")
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	assetsRepo := assetrepo.NewInMemoryAssetRepository()
	assetsService := assetapplication.NewAssetService(assetsRepo)

	assetClient := assets_client.NewAssetClient(assetsService)
	favRepo := favrepo.NewInMemoryFavouriteRepository()
	favService := favouritesapplication.NewFavouriteService(favRepo, assetClient)

	favHandler := httpapi.NewFavouritesHandler(
		favService,
		cfg.Pagination.DefaultPage,
		cfg.Pagination.DefaultLimit,
		cfg.Pagination.MaxLimit,
	)
	assetsHandler := httpapi.NewAssetsHandler(assetsService)

	pubKey, err := auth.LoadRSAPublicKey(cfg.Auth.PublicKeyPath)
	if err != nil {
		t.Fatalf("load public key: %v", err)
	}

	jwtCfg := auth.JWTConfig{
		PublicKey: pubKey,
		Issuer:    cfg.Auth.Issuer,
		Audience:  cfg.Auth.Audience,
	}

	router := httpapi.NewRouter(jwtCfg, favHandler, assetsHandler, false)
	ts := httptest.NewServer(router)

	token := mintToken(t, "user-1", cfg)
	return ts, token
}

func mintToken(t *testing.T, sub string, cfg config.AppConfig) string {
	t.Helper()

	keyBytes, err := os.ReadFile(cfg.Auth.PrivateKeyPath)
	if err != nil {
		t.Fatalf("read private key: %v", err)
	}
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		t.Fatalf("invalid private key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		t.Fatalf("parse private key: %v", err)
	}

	claims := jwt.RegisteredClaims{
		Subject:   sub,
		Issuer:    cfg.Auth.Issuer,
		Audience:  jwt.ClaimStrings{cfg.Auth.Audience},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(priv)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return ss
}

func doJSON(t *testing.T, client *http.Client, token string, method string, url string, body any, wantStatus int, out any) {
	t.Helper()

	var buf *bytes.Buffer
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		buf = bytes.NewBuffer(data)
	} else {
		buf = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequestWithContext(context.Background(), method, url, buf)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != wantStatus {
		t.Fatalf("%s %s status=%d want=%d", method, url, res.StatusCode, wantStatus)
	}

	if out != nil && res.ContentLength != 0 {
		if err := json.NewDecoder(res.Body).Decode(out); err != nil {
			t.Fatalf("decode response: %v", err)
		}
	}
}
