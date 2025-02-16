package main

import (
	"log"
	"platform-go-challenge/routes"
)

func main() {
	r := routes.SetupRouter()
	// Optional: Improve security by setting trusted proxies
	// Uncomment this line in production to avoid potential IP spoofing attacks
	// err := r.SetTrustedProxies([]string{"127.0.0.1"})
	// if err != nil {
	//     log.Fatalf("Failed to set trusted proxies: %v", err)
	// }
	// Start HTTPS server
	log.Println("🚀 Starting HTTPS server on https://localhost:8080")
	if err := r.RunTLS(":8080", "cert.pem", "key.pem"); err != nil {
		log.Fatalf("Failed to start HTTPS server: %v", err)
	}
}
