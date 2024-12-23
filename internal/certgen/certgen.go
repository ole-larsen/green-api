// certgen to generate cert
package certgen

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// GenerateCerts generates SAN certificates with no expiration and available for any network.
func GenerateCerts(dirName string) error {
	// Ensure the certs directory exists
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		return fmt.Errorf("failed to create certs directory: %v", err)
	}

	// 1. Generate the Private Key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	privateKeyFile, err := os.Create(filepath.Join(dirName, "server.key"))
	if err != nil {
		return fmt.Errorf("failed to create private key file: %v", err)
	}
	defer privateKeyFile.Close()

	err = pem.Encode(privateKeyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return fmt.Errorf("failed to write private key: %v", err)
	}

	// 2. Create the Certificate Signing Request (CSR)
	subject := pkix.Name{
		Country:            []string{"IN"},
		Province:           []string{"Goa"},
		Locality:           []string{"Panaji"},
		Organization:       []string{"Yandex"},
		OrganizationalUnit: []string{"IT"},
		CommonName:         "*.example.com", // Wildcard DNS for any subdomain
	}

	csrTemplate := x509.CertificateRequest{
		Subject: subject,
		DNSNames: []string{
			"*.example.com", // Wildcard domain for subdomains
		},
		IPAddresses: []net.IP{
			net.ParseIP("0.0.0.0"), // Allow localhost as a default IP address
		},
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create CSR: %v", err)
	}

	csrFile, err := os.Create(filepath.Join(dirName, "server.csr"))
	if err != nil {
		return fmt.Errorf("failed to create CSR file: %v", err)
	}
	defer csrFile.Close()

	err = pem.Encode(csrFile, &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to write CSR: %v", err)
	}

	// 3. Generate a Self-Signed Certificate (with no expiration)
	certTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(100 * 365 * 24 * time.Hour), // Valid for 100 years

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Add SAN extensions
	certTemplate.DNSNames = []string{"*.example.com"}           // Wildcard domain
	certTemplate.IPAddresses = []net.IP{net.ParseIP("0.0.0.0")} // Allow localhost as default

	// Self-sign the certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %v", err)
	}

	certFile, err := os.Create(filepath.Join(dirName, "server.crt"))
	if err != nil {
		return fmt.Errorf("failed to create certificate file: %v", err)
	}
	defer certFile.Close()

	err = pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to write certificate: %v", err)
	}

	fmt.Println("Certificates and key generated successfully in the " + dirName + "/ folder")

	return nil
}
