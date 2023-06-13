package main

import (
	"crypto/ed25519"
	"log"
)

func main() {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatal(err)
	}

	msg1 := []byte("The quick brown fox jumps over the lazy dog")

	msg2 := []byte("The quick brown fox jumps over the lazy dog")

	sig, err := priv.Sign(nil, msg1, &ed25519.Options{ //nolint:exhaustruct
		Context: "Example_ed25519ctx",
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := ed25519.VerifyWithOptions(pub, msg2, sig, &ed25519.Options{ //nolint:exhaustruct
		Context: "Example_ed25519ctx",
	}); err != nil {
		log.Fatal("invalid signature")
	}
}
