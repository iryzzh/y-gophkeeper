package tlsutil

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/iryzzh/gophkeeper/internal/rand"
)

// SecureDefault returns a tls.Config with reasonable, secure defaults set.
func SecureDefault() *tls.Config {
	return &tls.Config{
		// TLS 1.2 is the minimum we accept
		MinVersion: tls.VersionTLS12,
	}
}

// generateCert generates a new certificate.
func generateCert(commonName string) (certBlock, keyBlock *pem.Block, err error) {
	cert := &x509.Certificate{
		SerialNumber: new(big.Int).SetUint64(rand.Uint64()),
		Subject: pkix.Name{
			Organization: []string{"Gophkeeper"},
			Country:      []string{"RU"},
			CommonName:   commonName,
		},
		DNSNames:     []string{commonName},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	certBlock = &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}
	keyBlock = &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}

	return certBlock, keyBlock, nil
}

// NewCertificate generates and returns a new TLS certificate, saved to the given PEM files.
func NewCertificate(certFile, keyFile, commonName string) (tls.Certificate, error) {
	certBlock, keyBlock, err := generateCert(commonName)
	if err != nil {
		return tls.Certificate{}, err
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to save certificate: %w", err)
	}
	if err = pem.Encode(certOut, certBlock); err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to save certificate: %w", err)
	}
	if err = certOut.Close(); err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to save certificate: %w", err)
	}

	keyOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to save key: %w", err)
	}
	if err = pem.Encode(keyOut, keyBlock); err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to save key: %w", err)
	}
	if err = keyOut.Close(); err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to save key: %w", err)
	}

	return tls.X509KeyPair(pem.EncodeToMemory(certBlock), pem.EncodeToMemory(keyBlock))
}

// NewCertificateInMemory generates and returns a new TLS certificate, kept only in memory.
func NewCertificateInMemory(commonName string) (tls.Certificate, error) {
	certBlock, keyBlock, err := generateCert(commonName)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tls.X509KeyPair(pem.EncodeToMemory(certBlock), pem.EncodeToMemory(keyBlock))
}
