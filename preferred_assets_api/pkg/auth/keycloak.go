package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/cmd/api/config"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt"
)

type KeycloakClient struct {
	config   *config.KeycloakConfig
	client   *http.Client
	verifier *oidc.IDTokenVerifier
}

type CustomClaims struct {
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	PreferredName string `json:"preferred_username"`
	jwt.StandardClaims
}

func NewKeycloakClient(cfg *config.KeycloakConfig) (*KeycloakClient, error) {
	ctx := context.Background()

	URL := getKeycloakURL(cfg)
	log.Printf("%s", URL+"/realms/"+cfg.Realm)
	provider, err := oidc.NewProvider(ctx, URL+"/realms/"+cfg.Realm)
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID:          cfg.ClientID,
		SkipClientIDCheck: true,
	})

	return &KeycloakClient{
		config:   cfg,
		client:   &http.Client{Timeout: cfg.Timeout},
		verifier: verifier,
	}, nil
}

/*
	func NewKeycloakClient(cfg *config.KeycloakConfig) (*KeycloakClient, error) {
		ctx := context.Background()

		URL := getKeycloakURL(cfg)
		fullURL := URL + "/realms/" + cfg.Realm

		log.Printf("🔑 Initializing Keycloak Client:")
		log.Printf("   Base URL: %s", URL)
		log.Printf("   Full Realm URL: %s", fullURL)
		log.Printf("   Realm: %s", cfg.Realm)

		// Test connectivity first
		log.Printf("   Testing connectivity to Keycloak...")
		testClient := &http.Client{Timeout: 30 * time.Second}
		resp, err := testClient.Get(fullURL)
		if err != nil {
			log.Printf("   ❌ Connectivity test failed: %v", err)
			return nil, fmt.Errorf("cannot connect to Keycloak: %w", err)
		}
		defer resp.Body.Close()

		log.Printf("   ✅ Connectivity test passed - HTTP %d", resp.StatusCode)

		log.Printf("   Creating OIDC provider...")
		provider, err := oidc.NewProvider(ctx, fullURL)
		if err != nil {
			log.Printf("   ❌ Failed to create OIDC provider: %v", err)
			return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
		}
		log.Printf("   ✅ OIDC provider created successfully")

		verifier := provider.Verifier(&oidc.Config{
			ClientID:          cfg.ClientID,
			SkipClientIDCheck: true,
		})

		log.Printf("   ✅ Keycloak client initialized successfully")
		return &KeycloakClient{
			config:   cfg,
			client:   &http.Client{Timeout: cfg.Timeout},
			verifier: verifier,
		}, nil
	}

	func (k *KeycloakClient) testConfiguration(ctx context.Context, cfg *config.KeycloakConfig) error {
		// Try to get the OpenID configuration
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		URL := getKeycloakURL(cfg)
		println(URL)
		configURL := k.config.URL + "/realms/" + k.config.Realm + "/.well-known/openid-configuration"
		resp, err := k.client.Get(configURL)
		if err != nil {
			return fmt.Errorf("failed to fetch OpenID config: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("OpenID config returned status: %d", resp.StatusCode)
		}

		return nil
	}
*/
func (kc *KeycloakClient) VerifyToken(tokenString string) (*CustomClaims, error) {
	ctx := context.Background()
	if kc == nil {
		return nil, fmt.Errorf("keycloak client is not initialized")
	}
	if kc.verifier == nil {
		return nil, fmt.Errorf("verifier is not initialized - Keycloak client was not properly created")
	}
	idToken, err := kc.verifier.Verify(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %v", err)
	}

	var claims CustomClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %v", err)
	}

	return &claims, nil
}

func (kc *KeycloakClient) GetUserRoles(claims *CustomClaims) []string {
	var roles []string
	roles = append(roles, claims.RealmAccess.Roles...)
	if clientRoles, exists := claims.ResourceAccess[kc.config.ClientID]; exists {
		roles = append(roles, clientRoles.Roles...)
	}
	return roles
}

func getKeycloakURL(cfg *config.KeycloakConfig) string {
	// If running in Docker, use service name
	if os.Getenv("ENVIRONMENT") == "docker" {
		log.Println("Running in Docker environment")
		log.Println(cfg.URL)
		return cfg.URL
	}
	// If running locally, use localhost with mapped port
	log.Println("Running in local environment")
	return cfg.ExternalURL
}
