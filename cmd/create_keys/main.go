package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"
)

func main() {
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Printf("Cannot generate RSA key\n")
	}
	publickey := &privatekey.PublicKey

	// dump private key to file
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privatePem, err := os.Create("keys/private.pem")
	if err != nil {
		log.Panicf("error when create private.pem: %s\n", err)
	}
	err = pem.Encode(privatePem, privateKeyBlock)
	if err != nil {
		log.Panicf("error when encode private pem: %s\n", err)
	}

	// dump public key to file
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publickey)
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicPem, err := os.Create("keys/public.pem")
	if err != nil {
		log.Panicf("error when create public.pem: %s\n", err)
	}
	err = pem.Encode(publicPem, publicKeyBlock)
	if err != nil {
		log.Panicf("error when encode public pem: %s\n", err)
	}

	template := &x509.Certificate{
		Subject: pkix.Name{
			Country:            []string{"Mexico"},
			Organization:       []string{"Cultome"},
			CommonName:         "Cultome",
			Locality:           []string{"Coyoacan"},
			OrganizationalUnit: []string{"TV"},
		},
		SerialNumber:          big.NewInt(2019),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, template, publickey, privatekey)
	if err != nil {
		log.Panicf("Error creating certificate! %v\n", err)
	}

	pemfile, _ := os.Create("keys/cert.pem")
	pemkey := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	}

	pem.Encode(pemfile, pemkey)
	pemfile.Close()
}
