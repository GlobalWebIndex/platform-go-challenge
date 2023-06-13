package pki_test

import (
	"fmt"
	"log"
	"testing"

	"x-gwi/app/pki"
)

func ExampleNewPKI() {
	pk, err := pki.NewPKI()
	if err != nil {
		log.SetFlags(log.Lmicroseconds | log.Lshortfile)
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

func TestNewPKI(t *testing.T) {
	tests := []struct {
		name    string
		want    *pki.PKI
		wantErr bool
	}{
		// test cases.
		{name: "empty", want: &pki.PKI{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// got, err := pki.NewPKI()
			_, err := pki.NewPKI()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPKI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("NewPKI() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func BenchmarkNewPKI(b *testing.B) {
	b.Run("NewPKI", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pki.NewPKI()
		}
	})
}
