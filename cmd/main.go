// package main contains the entry point of the application.
package main

import (
	"fmt"
	"github.com/akazantzidis/gwi-ass/internal/app"
	"github.com/akazantzidis/gwi-ass/internal/infra"
	"github.com/akazantzidis/gwi-ass/internal/pkg/time"
	"github.com/akazantzidis/gwi-ass/internal/pkg/uuid"
	"log"
)

func main() {
	infraProviders := infra.NewInfraProviders()
	tp := time.NewTimeProvider()
	up := uuid.NewUUIDProvider()

	seedInitialUsers(infraProviders)

	appServices := app.NewServices(infraProviders.FavoriteRepository, infraProviders.NotificationService, infraProviders.UserRepository, infraProviders.RefreshTokenRepository, up, tp)

	infraHTTPServer := infra.NewHTTPServer(appServices)
	infraHTTPServer.ListenAndServe(":8080")
}

func seedInitialUsers(infra infra.Services) {
	user, err := infra.UserRepository.Add("alice", "password1", []string{"user"})
	if err != nil {
		log.Fatalf("failed to add user alice: %v", err)
	}
	fmt.Printf("alice %+v\n", user)

	user, err = infra.UserRepository.Add("bob", "password2", []string{"user", "admin"})
	if err != nil {
		log.Fatalf("failed to add user bob: %v", err)
	}
	fmt.Printf("bob: %+v\n", user)

	log.Println("Initial users seeded: alice, bob")
}
