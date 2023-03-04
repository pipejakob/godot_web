package godot_web

import (
	"crypto/x509"
	"testing"
)

func TestGenerateSelfSignedCertificate(t *testing.T) {
	cert, err := generateSelfSignedCertificate("192.168.0.100")

	if err != nil {
		t.Errorf("generateSelfSignedCertificate() returned unwanted error: %v", err)
	}
	if len(cert.Certificate) != 1 {
		t.Errorf("generateSelfSignedCertificate() wanted Certificate length 1, got %v",
			len(cert.Certificate))
	}
	parsed, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		t.Errorf("generateSelfSignedCertificate() returned unparseable certificate: %v", err)
	}
	expectedCommon := "192.168.0.100"
	if parsed.Subject.CommonName != expectedCommon {
		t.Errorf("generateSelfSignedCertificate() wanted CommonName %q, got %q",
			expectedCommon, parsed.Subject.CommonName)
	}
	if !parsed.IsCA {
		t.Errorf("generateSelfSignedCertificate() wanted IsCA = true, but was false")
	}
}
