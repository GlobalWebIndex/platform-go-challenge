package pki

import (
	"crypto/ed25519"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type ED25519 struct {
	key     ed25519.PrivateKey // x509
	pub     ed25519.PublicKey
	certTLS *tls.Certificate // tls
	keyDER  []byte
	keyPEM  []byte // x509 tls
	pubDER  []byte
	pubPEM  []byte
	certDER []byte
	certPEM []byte // x509 tls
}

func newED25519() (*ED25519, error) {
	// ED25519 not yet supported by browsers
	_, key, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, fmt.Errorf("ed25519.GenerateKey: %w", err)
	}

	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("x509.MarshalPKCS8PrivateKey: %w", err)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Headers: map[string]string{}, Bytes: keyDER})

	pub, _ := key.Public().(ed25519.PublicKey)

	pubDER, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, fmt.Errorf("x509.MarshalPKCS8PrivateKey: %w", err)
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Headers: map[string]string{}, Bytes: pubDER})

	return &ED25519{key: key, keyDER: keyDER, keyPEM: keyPEM, //nolint:exhaustruct
		pub: pub, pubDER: pubDER, pubPEM: pubPEM,
		// certDER: []byte{}, certPEM: []byte{},
	}, nil
}

// certPEM, keyPEM
func (c *ED25519) keyPairPEM() ([]byte, []byte) {
	return c.certPEM, c.keyPEM
}

func (c *ED25519) tlsX509KeyPairCert() error {
	certTLS, err := tls.X509KeyPair(c.keyPairPEM())
	if err != nil {
		return fmt.Errorf("tls.X509KeyPair: %w", err)
	}

	c.certTLS = &certTLS

	return nil
}
