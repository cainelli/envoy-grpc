package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/go-jose/go-jose"
)

const JWKSPath = "/.well-known/jwks.json"

type JWTAuthn struct {
	KeySet *jose.JSONWebKeySet
}

func (j *JWTAuthn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != JWKSPath {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if j.KeySet == nil {
		jwks, err := NewKeySet()
		if err != nil {
			http.Error(w, "error creating JWKS", http.StatusInternalServerError)
			return
		}
		j.KeySet = jwks
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(j.KeySet); err != nil {
		http.Error(w, "error encoding JWKS", http.StatusInternalServerError)
		return
	}
}

func NewKeySet() (*jose.JSONWebKeySet, error) {
	rsaKey, _ := rsa.GenerateKey(rand.Reader, 2048)           // XXX Check err
	serialNumber, _ := rand.Int(rand.Reader, big.NewInt(100)) // XXX Check err

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Example Co"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(2 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &rsaKey.PublicKey, rsaKey)
	if err != nil {
		return nil, fmt.Errorf("error creating certificate: %w", err)
	}

	certificate, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %w", err)
	}

	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{{
			Certificates: []*x509.Certificate{certificate},
			Key:          &rsaKey.PublicKey,
			KeyID:        "1",
			Use:          "sig",
		}},
	}, nil
}
