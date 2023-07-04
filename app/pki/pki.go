package pki

import (
	"crypto"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/fs"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PKI struct {
	selfCA  *ED25519 // crt
	ed25519 *ED25519 // tls
	ecdsa   *ECDSA   // tls jwt
	rsa     *RSA     // tls jwt encr
	roots   *x509.CertPool
}

type casePKI string

const (
	_selfCA  casePKI = "selfCA"
	_ed25519 casePKI = "ed25519"
	_ecdsa   casePKI = "ecdsa"
	_rsa     casePKI = "rsa"
)

func NewPKI() (*PKI, error) {
	sc, err := newED25519()
	if err != nil {
		return nil, fmt.Errorf("newED25519: %w", err)
	}

	ed, err := newED25519()
	if err != nil {
		return nil, fmt.Errorf("newED25519: %w", err)
	}

	ec, err := newECDSA()
	if err != nil {
		return nil, fmt.Errorf("newECDSA: %w", err)
	}

	rs, err := newRSA()
	if err != nil {
		return nil, fmt.Errorf("newRSA: %w", err)
	}

	// https://pkg.go.dev/crypto/x509#example-Certificate.Verify (1/4)
	// SSL_CERT_FILE and SSL_CERT_DIR can be used to override the system default locations
	//
	roots, err := x509.SystemCertPool() // x509.NewCertPool()  x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("x509.SystemCertPool: %w", err)
	}

	p := &PKI{selfCA: sc, ed25519: ed, ecdsa: ec, rsa: rs, roots: roots}

	err = p.selfCertify()
	if err != nil {
		return nil, fmt.Errorf("p.selfCertify: %w", err)
	}

	return p, nil
}

func (p *PKI) selfCertify() error { //nolint:funlen,cyclop
	notBefore := time.Now().UTC() // .Add(-366 * 24 * time.Hour) // .
	notAfter := notBefore.Add(365 * 24 * time.Hour)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128) //nolint:gomnd

	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("rand.Int: %w", err)
	}

	// selfCA
	templCA := x509.Certificate{ //nolint:exhaustruct
		SerialNumber: serialNumber,
		Subject: pkix.Name{ //nolint:exhaustruct
			Organization: []string{"Acme Co"},
			CommonName:   "Self Root CA ed25519",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	var (
		IPAddresses []net.IP
		DNSNames    []string
	)

	// env.Env("SERVER_GRPC_ADDRESS", ":9090")
	host := "127.0.0.1,[::],localhost,localhost.localdomain"
	hosts := strings.Split(host, ",")

	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			IPAddresses = append(IPAddresses, ip)
		} else {
			DNSNames = append(DNSNames, h)
		}
	}

	// ED25519
	serialNumber, _ = rand.Int(rand.Reader, serialNumberLimit)
	templED := templCA
	templED.IsCA = false
	templED.SerialNumber = serialNumber
	templED.Subject.CommonName = "ed25519_cert"
	templED.KeyUsage = x509.KeyUsageDigitalSignature
	templED.DNSNames = DNSNames
	templED.IPAddresses = IPAddresses

	// ECDSA
	serialNumber, _ = rand.Int(rand.Reader, serialNumberLimit)
	templEC := templED
	templED.SerialNumber = serialNumber
	templED.Subject.CommonName = "ecdsa_cert"

	// RSA
	serialNumber, _ = rand.Int(rand.Reader, serialNumberLimit)
	templRS := templED
	templED.SerialNumber = serialNumber
	templED.Subject.CommonName = "rsa_cert"
	templRS.KeyUsage = x509.KeyUsageDigitalSignature |
		x509.KeyUsageKeyEncipherment |
		x509.KeyUsageDataEncipherment

	// if templED.IsCA || templEC.IsCA || templRS.IsCA {}

	// cert self CA by self CA
	p.selfCA.certDER, err = x509.CreateCertificate(rand.Reader, &templCA, &templCA, p.selfCA.pub, p.keyCA())
	if err != nil {
		return fmt.Errorf("x509.CreateCertificate: %w", err)
	}

	p.selfCA.certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Headers: map[string]string{},
		Bytes: p.selfCA.certDER})

	// cert ED25519 by self CA
	p.ed25519.certDER, err = x509.CreateCertificate(rand.Reader, &templED, &templCA, p.ed25519.pub, p.keyCA())
	if err != nil {
		return fmt.Errorf("x509.CreateCertificate: %w", err)
	}

	p.ed25519.certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Headers: map[string]string{},
		Bytes: p.ed25519.certDER})

	// cert ECDSA by self CA
	p.ecdsa.certDER, err = x509.CreateCertificate(rand.Reader, &templEC, &templCA, p.ecdsa.pub, p.keyCA())
	if err != nil {
		return fmt.Errorf("x509.CreateCertificate: %w", err)
	}

	p.ecdsa.certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Headers: map[string]string{},
		Bytes: p.ecdsa.certDER})

	// x509.ParseCertificate(p.ecdsa.certDER)

	// cert RSA by self CA
	p.rsa.certDER, err = x509.CreateCertificate(rand.Reader, &templRS, &templCA, p.rsa.pub, p.keyCA())
	if err != nil {
		return fmt.Errorf("x509.CreateCertificate: %w", err)
	}

	p.rsa.certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Headers: map[string]string{},
		Bytes: p.rsa.certDER})

	// https://pkg.go.dev/crypto/x509#example-Certificate.Verify (2/4)
	//
	if ok := p.roots.AppendCertsFromPEM(p.selfCA.certPEM); !ok {
		return fmt.Errorf("roots.AppendCertsFromPEM(p.selfCA.certPEM): !ok") //nolint:goerr113
	}
	//
	// Verify all self signed certificates
	if err = p.VerifyCertPEM(); err != nil {
		return fmt.Errorf("p.Validate: %w", err)
	}
	//
	// validation of all self signed certificates
	if err = p.ValidateTLSX509KeyPairCert(); err != nil {
		return fmt.Errorf("p.ValidateTLSX509KeyPairCert: %w", err)
	}
	//
	// saving of all self signed certificates and keys
	if err = p.saveKeyPairCerts(); err != nil {
		return fmt.Errorf("p.saveKeyPairCerts(): %w", err)
	}

	return nil
}

func (p *PKI) VerifyCertPEM() error {
	// https://pkg.go.dev/crypto/x509#example-Certificate.Verify (3/4)
	//
	if err := p.verifyCertPEM(p.selfCA.certPEM); err != nil {
		return fmt.Errorf("p.verifyCertPEM selfCA: %w", err)
	}

	if err := p.verifyCertPEM(p.ed25519.certPEM); err != nil {
		return fmt.Errorf("p.verifyCertPEM ed25519: %w", err)
	}

	if err := p.verifyCertPEM(p.ecdsa.certPEM); err != nil {
		return fmt.Errorf("p.verifyCertPEM ecdsa: %w", err)
	}

	if err := p.verifyCertPEM(p.rsa.certPEM); err != nil {
		return fmt.Errorf("p.verifyCertPEM rsa: %w", err)
	}

	return nil
}

func (p *PKI) verifyCertPEM(certPEM []byte) error {
	// https://pkg.go.dev/crypto/x509#example-Certificate.Verify (4/4)
	//
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return fmt.Errorf("pem.Decode: block is nil") //nolint:goerr113
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("x509.ParseCertificate: %w", err)
	}

	opts := x509.VerifyOptions{ //nolint:exhaustruct
		// DNSName: "",
		Roots: p.roots,
	}

	if _, err := cert.Verify(opts); err != nil {
		return fmt.Errorf("cert.Verify: %w", err)
	}

	return nil
}

func (p *PKI) ValidateTLSX509KeyPairCert() error {
	if err := p.selfCA.tlsX509KeyPairCert(); err != nil {
		return fmt.Errorf("selfCA.tlsX509KeyPairCert: %w", err)
	}

	if err := p.ed25519.tlsX509KeyPairCert(); err != nil {
		return fmt.Errorf("ed25519.tlsX509KeyPairCert: %w", err)
	}

	if err := p.ecdsa.tlsX509KeyPairCert(); err != nil {
		return fmt.Errorf("ecdsa.tlsX509KeyPairCert: %w", err)
	}

	if err := p.rsa.tlsX509KeyPairCert(); err != nil {
		return fmt.Errorf("rsa.tlsX509KeyPairCert: %w", err)
	}

	return nil
}

func (p *PKI) saveKeyPairCerts() error {
	// switch casePKI {
	// case _selfCA:
	// case _ed25519:
	// case _ecdsa:
	// case _rsa:
	// default:
	// 	return fmt.Errorf("pki: casePKI: %v", casePKI) //nolint:goerr113
	// }
	if err := p.saveKeyPairCert(_selfCA, p.selfCA.certPEM, p.selfCA.keyPEM); err != nil {
		return fmt.Errorf("pki: p.saveKeyPairCert: selfCA %w", err)
	}

	if err := p.saveKeyPairCert(_ed25519, p.ed25519.certPEM, p.ed25519.keyPEM); err != nil {
		return fmt.Errorf("pki: p.saveKeyPairCert: ed25519 %w", err)
	}

	if err := p.saveKeyPairCert(_ecdsa, p.ecdsa.certPEM, p.ecdsa.keyPEM); err != nil {
		return fmt.Errorf("pki: p.saveKeyPairCert: ecdsa %w", err)
	}

	if err := p.saveKeyPairCert(_rsa, p.rsa.certPEM, p.rsa.keyPEM); err != nil {
		return fmt.Errorf("pki: p.saveKeyPairCert: rsa %w", err)
	}

	return nil
}

func (p *PKI) saveKeyPairCert(casePKI casePKI, certPEMBlock []byte, keyPEMBlock []byte) error {
	const (
		permDir fs.FileMode = 0740 // 740 rwxr-----
		permCer fs.FileMode = 0640 // 640 rw-r-----
		permKey fs.FileMode = 0600 // 600 rw-------
	)

	saveDir := filepath.Join("test", "certs", "auto-self")
	saveCer := filepath.Join(saveDir, fmt.Sprintf("%s-cert.pem", casePKI))
	saveKey := filepath.Join(saveDir, fmt.Sprintf("%s-key.pem", casePKI))

	// if err := os.RemoveAll(filepath.Join("test", "certs")); err != nil {
	// 	return fmt.Errorf("pki: os.RemoveAll %w", err)
	// }

	if err := os.MkdirAll(saveDir, permDir); err != nil { //  && !os.IsExist(err)
		return fmt.Errorf("pki: os.MkdirAll %w", err)
	}

	if err := os.WriteFile(saveCer, certPEMBlock, permCer); err != nil {
		return fmt.Errorf("pki: os.WriteFile cert %w", err)
	}

	if err := os.WriteFile(saveKey, keyPEMBlock, permKey); err != nil {
		return fmt.Errorf("pki: os.WriteFile key %w", err)
	}

	if err := os.Chmod(saveKey, permKey); err != nil {
		return fmt.Errorf("pki: os.Chmod key %w", err)
	}

	return nil
}

// keyCA use for self signing certificates
func (p *PKI) keyCA() crypto.PrivateKey {
	return p.selfCA.key
}

// CACertPEM use by tls client cert
func (p *PKI) CACertPEM() []byte {
	return p.selfCA.certPEM
}

// CertPool()
func (p *PKI) CertPool() *x509.CertPool {
	return p.roots
}

func (p *PKI) TLSConfigServerGRPC() *tls.Config {
	return &tls.Config{ //nolint:exhaustruct
		Certificates: []tls.Certificate{*p.ed25519.certTLS}, // web ecdsa secp384r1 (P-384); ed25519 not supported in browser
		MinVersion:   tls.VersionTLS13,
		// CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP521, tls.CurveP384, tls.CurveP256},
		ClientCAs:  p.CertPool(),
		ClientAuth: tls.VerifyClientCertIfGiven, // NoClientCert VerifyClientCertIfGiven RequireAndVerifyClientCert
	}
}

func (p *PKI) TLSConfigServerRESTGW() *tls.Config {
	return &tls.Config{ //nolint:exhaustruct
		Certificates: []tls.Certificate{*p.ecdsa.certTLS}, // use ecdsa secp384r1 (P-384); ed25519 not supported in browser
		MinVersion:   tls.VersionTLS13,
		// CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP521, tls.CurveP384, tls.CurveP256},
		ClientCAs:  p.CertPool(),
		ClientAuth: tls.VerifyClientCertIfGiven, // NoClientCert VerifyClientCertIfGiven RequireAndVerifyClientCert
	}
}

func (p *PKI) TLSConfigDial() *tls.Config {
	return &tls.Config{ //nolint:exhaustruct
		// Certificates:       []tls.Certificate{certTLSdial},
		MinVersion: tls.VersionTLS13,
		RootCAs:    p.CertPool(),
		// InsecureSkipVerify: false,
	}
}
