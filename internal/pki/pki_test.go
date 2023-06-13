package pki_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"
)

/*
import "x-gwi/internal/pki"
func ExampleNewPKI() {
	pk, err := pki.NewPKI()
	if err != nil {
		log.Fatal(err)
		// fmt.Printf("pki.NewPKI() pki: %+v err: %+v\n", pk, err)
	} else {
		fmt.Printf("pki.NewPKI OK\n")
	}

	fmt.Printf("VerifyCertPEM err: %+v\n", pk.VerifyCertPEM())
	fmt.Printf("ValidateTLSX509KeyPairCert err: %+v\n", pk.ValidateTLSX509KeyPairCert())
	// fmt.Printf("pkiPub: %+v\npkiCert: %+v\npkiCACert: %+v\n", string(pk.PublicKeyPEM()), string(pk.PublicKeyCertPEM()), string(pk.CACertPEM()))
	// fmt.Printf("pkiTLSConfig: %+v\n\n", len(fmt.Sprintf("%+v", *pk.TLSConfig())))

	// fmt.Printf("caPub: %+v\ncaCert: %+v\n", string(pk.selfCA.PublicKeyPEM()), string(pk.selfCA.PublicKeyCertPEM()))
	// fmt.Printf("edPub: %+v\nedCert: %+v\n", string(pk.ed25519.PublicKeyPEM()), string(pk.ed25519.PublicKeyCertPEM()))
	// fmt.Printf("ecPub: %+v\necCert: %+v\n", string(pk.ecdsa.PublicKeyPEM()), string(pk.ecdsa.PublicKeyCertPEM()))
	// fmt.Printf("rsPub: %+v\nrsCert: %+v\n", string(pk.RSA.PublicKeyPEM()), string(pk.RSA.PublicKeyCertPEM()))

	// Output:
	//
	// pki.NewPKI OK
	// VerifyCertPEM err: <nil>
	// ValidateTLSX509KeyPairCert err: <nil>
}
*/

// ECDSA, ED25519 and RSA
func Example() {
	// ED25519
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("pub: %T: %+v\npub: %T: %+v\nprv: %T: %+v\n", pub, pub, priv.Public(), priv.Public(), priv, priv)
	fmt.Printf("pub: %d, prv: %d\n", len(pub), len(priv))

	msg := []byte("The quick brown fox jumps over the lazy dog")

	msg2 := []byte("The quick brown fox jumps over the lazy dog")

	sig, err := priv.Sign(nil, msg, &ed25519.Options{
		Context: "Example_ed25519ctx",
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := ed25519.VerifyWithOptions(pub, msg2, sig, &ed25519.Options{
		Context: "Example_ed25519ctx",
	}); err != nil {
		log.Fatal("invalid signature")
	}

	// RSA
	// Generate a 2048-bits key
	privateKey, publicKey := generateKeyPairRSA(4096) // 2048 3072 4096 5120 6144 7168 8192 9216
	message := []byte("super secret message")
	t := time.Now()
	cipherText, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, message, nil)
	t21 := time.Since(t)
	fmt.Printf("%v Encrypted: %v\n", t21, cipherText)
	t = time.Now()
	decMessage, _ := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, cipherText, nil)
	t21 = time.Since(t)
	fmt.Printf("%v Decrypted: %s\n", t21, string(decMessage))
	fmt.Printf("pub: %d, prv: %d\n", publicKey.Size(), privateKey.Size())

	// 	ECDSA Keys
	// Most browsers support  secp256r1 (P-256)  and  secp384r1 (P-384)
	// chrome/edge err on (P-521)

	// Output:
	//
	//
}

func generateKeyPairRSA(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	// This method requires a random number of bits.
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// The public key is part of the PrivateKey struct
	return privateKey, &privateKey.PublicKey
}

func ExampleX509KeyPair_httpServer() {
	certPem := []byte(`-----BEGIN CERTIFICATE-----
MIIBhTCCASugAwIBAgIQIRi6zePL6mKjOipn+dNuaTAKBggqhkjOPQQDAjASMRAw
DgYDVQQKEwdBY21lIENvMB4XDTE3MTAyMDE5NDMwNloXDTE4MTAyMDE5NDMwNlow
EjEQMA4GA1UEChMHQWNtZSBDbzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABD0d
7VNhbWvZLWPuj/RtHFjvtJBEwOkhbN/BnnE8rnZR8+sbwnc/KhCk3FhnpHZnQz7B
5aETbbIgmuvewdjvSBSjYzBhMA4GA1UdDwEB/wQEAwICpDATBgNVHSUEDDAKBggr
BgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1UdEQQiMCCCDmxvY2FsaG9zdDo1
NDUzgg4xMjcuMC4wLjE6NTQ1MzAKBggqhkjOPQQDAgNIADBFAiEA2zpJEPQyz6/l
Wf86aX6PepsntZv2GYlA5UpabfT2EZICICpJ5h/iI+i341gBmLiAFQOyTDT+/wQc
6MF9+Yw1Yy0t
-----END CERTIFICATE-----`)
	keyPem := []byte(`-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIIrYSSNQFaA2Hwf1duRSxKtLYX5CB04fSeQ6tF1aY/PuoAoGCCqGSM49
AwEHoUQDQgAEPR3tU2Fta9ktY+6P9G0cWO+0kETA6SFs38GecTyudlHz6xvCdz8q
EKTcWGekdmdDPsHloRNtsiCa697B2O9IFA==
-----END EC PRIVATE KEY-----`)
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		log.Fatal(err)
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	srv := &http.Server{
		TLSConfig:    cfg,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	log.Fatal(srv.ListenAndServeTLS("", ""))

	// Output:
	//
	//
}
