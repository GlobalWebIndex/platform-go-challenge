package pki

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type ECDSA struct {
	key     *ecdsa.PrivateKey // x509
	pub     *ecdsa.PublicKey
	certTLS *tls.Certificate // tls
	keyDER  []byte
	keyPEM  []byte // tls
	pubDER  []byte
	pubPEM  []byte
	certDER []byte
	certPEM []byte // tls
}

func newECDSA() (*ECDSA, error) {
	// ECDSA Keys Most browsers support  secp256r1 (P-256)  and  secp384r1 (P-384)
	// chrome/edge err on (P-521)
	curve := elliptic.P384() // P224() P256() P384() P521()

	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("ecdsa.GenerateKey: %w", err)
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

	return &ECDSA{key: key, keyDER: keyDER, keyPEM: keyPEM, //nolint:exhaustruct
		pub: pub, pubDER: pubDER, pubPEM: pubPEM,
		// certDER: []byte{}, certPEM: []byte{}, certTLS: tls.Certificate{},
	}, nil
}

// certPEM, keyPEM
func (c *ECDSA) keyPairPEM() ([]byte, []byte) {
	return c.certPEM, c.keyPEM
}

func (c *ECDSA) tlsX509KeyPairCert() error {
	certTLS, err := tls.X509KeyPair(c.keyPairPEM())
	if err != nil {
		return fmt.Errorf("tls.X509KeyPair: %w", err)
	}

	c.certTLS = &certTLS

	return nil
}
