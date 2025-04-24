package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"runtime"
	"strings"
	"time"

	"github.com/briskt/go-htmx-app/action"
	"github.com/briskt/go-htmx-app/app"
	"github.com/briskt/go-htmx-app/log"
)

func main() {
	log.Init(
		log.UseCommit(strings.TrimSpace(app.Commit)),
		log.UseEnv(app.Env.AppEnv),
		log.UseLevel(app.Env.LogLevel),
		log.UsePretty(app.Env.AppEnv == app.EnvDevelopment),
		log.UseRemote(app.Env.AppEnv != app.EnvTest),
	)

	log.WithFields(log.Fields{
		"appEnv":     app.Env.AppName,
		"goVersion":  runtime.Version(),
		"commitHash": strings.TrimSpace(app.Commit),
	}).Info("app starting")

	db, err := app.OpenDatabase()
	if err != nil {
		log.Fatalf("database error: %s", err)
	}

	emailService, err := app.NewEmailService()
	if err != nil {
		log.Fatalf("error creating email service: %s", err)
	}

	a := action.NewApp(&action.Config{
		DB:           db,
		EmailService: emailService,
	})

	if app.Env.DisableTLS {
		log.Fatal(a.Start(":80"))
	} else {
		cert, key, err := generateCert()
		if err != nil {
			log.Fatalf("failed to generate cert: %s", err)
		}

		log.Fatal(a.StartTLS(":443", cert, key))
	}
}

// This code was informed by crypto/tls/generate_cert.go in the Go source repository

// generateCert creates a new self-signed certificate
func generateCert() (cert []byte, key []byte, err error) {
	var privateKey *rsa.PrivateKey
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"SIL International"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}

	// for self-cert
	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	// create cert
	var certBuffer bytes.Buffer
	if err = pem.Encode(&certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return nil, nil, fmt.Errorf("failed to write cert: %w", err)
	}
	cert = certBuffer.Bytes()

	// create key
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal private key: %w", err)
	}
	var keyBuffer bytes.Buffer
	if err = pem.Encode(&keyBuffer, &pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyBytes}); err != nil {
		return nil, nil, fmt.Errorf("failed to write key: %w", err)
	}
	key = keyBuffer.Bytes()

	return
}
