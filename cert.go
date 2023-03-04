package godot_web

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"time"
)

func generateSelfSignedCertificate(commonName string) (tls.Certificate, error) {
	now := time.Now()
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(now.Unix()),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"godot_web"},
		},
		NotBefore:             now,
		NotAfter:              now.AddDate(1, 0, 0), // +1 year
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		KeyUsage: x509.KeyUsageCertSign |
			x509.KeyUsageDigitalSignature |
			x509.KeyUsageKeyEncipherment,
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("error generating RSA key: %v", err)
	}

	signed, err := x509.CreateCertificate(rand.Reader, cert, cert, key.Public(), key)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("error signing certificate: %v", err)
	}

	return tls.Certificate{
		Certificate: [][]byte{
			signed,
		},
		PrivateKey: key,
	}, nil
}
