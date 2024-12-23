package certgen_test

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/ole-larsen/green-api/internal/certgen"
)

// TestGenerateCerts checks if the GenerateCerts function creates the expected certificate files.
func TestGenerateCerts(t *testing.T) {
	// Create a temporary directory for the certs
	dir := t.TempDir()

	// Call the GenerateCerts function
	err := certgen.GenerateCerts(dir)
	if err != nil {
		t.Fatalf("GenerateCerts failed: %v", err)
	}

	// Check if server.key exists
	keyPath := filepath.Join(dir, "server.key")
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		t.Errorf("server.key does not exist in %s", dir)
	}

	// Check if server.crt exists
	crtPath := filepath.Join(dir, "server.crt")
	if _, err := os.Stat(crtPath); os.IsNotExist(err) {
		t.Errorf("server.crt does not exist in %s", dir)
	}

	// Check if server.csr exists
	csrPath := filepath.Join(dir, "server.csr")
	if _, err := os.Stat(csrPath); os.IsNotExist(err) {
		t.Errorf("server.csr does not exist in %s", dir)
	}

	// Validate the server.crt as a valid PEM-encoded certificate
	validatePEMCertificate(t, crtPath)

	// Validate the server.key as a valid PEM-encoded RSA private key
	validatePEMKey(t, keyPath)
}

// validatePEMCertificate validates that the given file is a PEM-encoded certificate.
func validatePEMCertificate(t *testing.T, certPath string) {
	// Read the certificate file
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		t.Fatalf("failed to read certificate file %s: %v", certPath, err)
	}

	// Decode the PEM block
	block, _ := pem.Decode(certBytes)
	if block == nil {
		t.Fatalf("failed to decode PEM block from certificate")
	}

	// Parse the x509 certificate
	_, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("failed to parse certificate: %v", err)
	}
}

// validatePEMKey validates that the given file is a PEM-encoded private key.
func validatePEMKey(t *testing.T, keyPath string) {
	// Read the private key file
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		t.Fatalf("failed to read private key file %s: %v", keyPath, err)
	}

	// Decode the PEM block
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		t.Fatalf("failed to decode PEM block from private key")
	}

	// Parse the RSA private key
	_, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		t.Fatalf("failed to parse private key: %v", err)
	}
}
