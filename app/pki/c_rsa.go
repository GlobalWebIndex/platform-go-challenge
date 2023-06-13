package pki

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type RSA struct {
	key     *rsa.PrivateKey // x509
	pub     *rsa.PublicKey
	certTLS *tls.Certificate // tls
	keyDER  []byte
	keyPEM  []byte // tls
	pubDER  []byte
	pubPEM  []byte
	certDER []byte
	certPEM []byte // tls
}

func newRSA() (*RSA, error) {
	const bits int = 4096 // 2048 3072 4096 5120 6144 7168 8192 9216

	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("rsa.GenerateKey %w", err)
	}

	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("x509.MarshalPKCS8PrivateKey: %w", err)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Headers: map[string]string{}, Bytes: keyDER})

	pub := &key.PublicKey

	pubDER, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, fmt.Errorf("x509.MarshalPKCS8PrivateKey: %w", err)
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Headers: map[string]string{}, Bytes: pubDER})

	return &RSA{key: key, keyDER: keyDER, keyPEM: keyPEM, //nolint:exhaustruct
		pub: pub, pubDER: pubDER, pubPEM: pubPEM,
		// certDER: []byte{}, certPEM: []byte{},
	}, nil
}

// certPEM, keyPEM
func (c *RSA) keyPairPEM() ([]byte, []byte) {
	return c.certPEM, c.keyPEM
}

func (c *RSA) tlsX509KeyPairCert() error {
	certTLS, err := tls.X509KeyPair(c.keyPairPEM())
	if err != nil {
		return fmt.Errorf("tls.X509KeyPair: %w", err)
	}

	c.certTLS = &certTLS

	return nil
}
