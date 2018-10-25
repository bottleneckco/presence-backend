package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"strings"
)

// Generate an RSA private key in dotenv format
func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Panic(err)
	}

	pemPrivateBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyPEM := string(pem.EncodeToMemory(pemPrivateBlock))

	fmt.Printf("JWT_PRIVATE_KEY=\"%s\"\n", strings.Replace(privateKeyPEM, "\n", "\\n", -1))
}
