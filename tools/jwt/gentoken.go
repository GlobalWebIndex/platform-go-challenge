package main

import (
	"fmt"
	"log"
	"os"
	"platform-go-challenge/internal/api/middleware"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(".env.test"); err != nil {
		log.Fatal("Error loading .env.test file")
	}

	user := os.Args[1]
	if len(user) == 0 {
		fmt.Println("No user flag detected")
		fmt.Println("defaulting to user1")
		user = "user1"
	}
	// Generate a token for test user
	token, err := middleware.GenerateToken(user)
	if err != nil {
		log.Fatal("Error generating token:", err)
	}

	fmt.Printf("\nJWT_TOKEN=%s\n", token)
}
