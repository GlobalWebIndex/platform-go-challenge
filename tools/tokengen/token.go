package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	userID := flag.String("user", "", "user ID to embed in the token (required)")
	secret := flag.String("secret", "", "HMAC signing secret (or set JWT_SECRET env var)")
	expiry := flag.Duration("exp", 24*time.Hour, "token expiry duration (e.g. 1h, 72h)")
	flag.Parse()

	if *userID == "" {
		fmt.Fprintln(os.Stderr, "error: -user flag is required")
		flag.Usage()
		os.Exit(1)
	}

	signingSecret := *secret
	if signingSecret == "" {
		signingSecret = os.Getenv("JWT_SECRET")
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub": *userID,
		"iat": now.Unix(),
		"exp": now.Add(*expiry).Unix(),
	}

	var signed string
	if signingSecret == "" {
		token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		var err error
		signed, err = token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating token: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "Warning: token is unsigned (alg=none); do not use in production")
	} else {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		var err error
		signed, err = token.SignedString([]byte(signingSecret))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error signing token: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Fprintf(os.Stderr, "Token for user %s (expires %s):\n", *userID, now.Add(*expiry).Format(time.RFC3339))
	fmt.Println(signed)
}
