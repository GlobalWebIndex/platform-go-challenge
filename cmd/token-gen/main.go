package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/manos/favourites/config"
)

func main() {
	var (
		env      = flag.String("env", "", "Config environment (local/dev/test/prod). Overrides APP_ENV.")
		user     = flag.String("user", "user-1", "User id to set as JWT subject (sub)")
		issuer   = flag.String("iss", "", "JWT issuer (iss)")
		audience = flag.String("aud", "", "JWT audience (aud)")
		minutes  = flag.Int("minutes", 15, "Token validity in minutes")
		keyPath  = flag.String("key", "", "Path to RSA private key PEM (for signing)")
	)

	flag.Parse()

	cfg, err := config.LoadWithEnv(*env)
	if err != nil {
		panic(fmt.Sprintf("load config: %v", err))
	}

	if *issuer == "" {
		*issuer = cfg.Auth.Issuer
	}
	if *audience == "" {
		*audience = cfg.Auth.Audience
	}
	if *keyPath == "" {
		*keyPath = cfg.Auth.PrivateKeyPath
	}

	fmt.Println("Generating JWT token...", *user, *issuer, *audience, *minutes, *keyPath)

	keyData, err := os.ReadFile(*keyPath)
	if err != nil {
		panic(err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		panic(err)
	}

	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   *user, // ✅ critical: user identity goes here
		Issuer:    *issuer,
		Audience:  jwt.ClaimStrings{*audience},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(*minutes) * time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(privKey)
	if err != nil {
		panic(err)
	}

	fmt.Println(signed)
}
