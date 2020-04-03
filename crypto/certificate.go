package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
	"net"
	"time"
)

// Some Operation related to Public & Private Key
// Generated CSR using algorithm rsa:2048
// example:
//      openssl req -new -nodes -newkey rsa:2048 \ -keyout msslclient.key -out msslclient.csr -subj \
//      "/C=ID/OU=mtf.comapi-client/O=ICCP/CN=Tokopedia-ICCP"

type CACertificate struct {
	PrivateKey  []byte
	CSR         []byte
	Certificate []byte
}

// Generate Certificate Signing Request (CSR) and private key
func GenerateCA(csrTemplate *x509.CertificateRequest, certTemplate *x509.Certificate, parentTemplate *x509.Certificate) (*CACertificate, error) {
	// 1. Create PrivateKey
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, errors.New("GenerateCA " + err.Error())
	}

	// Encode privateKey
	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	// 2. Create CSR template
	if csrTemplate == nil {
		csrTemplate = &x509.CertificateRequest{
			Subject: pkix.Name{
				Organization:  []string{"ESEMKA"},
				Country:       []string{"ID"},
				Province:      []string{""},
				Locality:      []string{"Jakarta"},
				StreetAddress: []string{"Jl Guru Mughni"},
				PostalCode:    []string{"54393"},
			},
		}
	}
	// create CSR
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, caPrivKey)
	if err != nil {
		return nil, errors.New("GenerateCA " + err.Error())
	}

	// Encode CSR
	csrPEM := new(bytes.Buffer)
	pem.Encode(csrPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: csrBytes,
	})

	// 3. Create certificate Template
	if certTemplate == nil {
		certTemplate = &x509.Certificate{
			SerialNumber: big.NewInt(2019),
			Subject: pkix.Name{
				Organization:  []string{"ESTEEM"},
				Country:       []string{"Indo"},
				Province:      []string{"DKI"},
				Locality:      []string{"Jaksel"},
				StreetAddress: []string{"Jl Guru Mughni"},
				PostalCode:    []string{"94016"},
			},
			IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
			NotBefore:    time.Now(),
			NotAfter:     time.Now().AddDate(10, 0, 0),
			SubjectKeyId: []byte{1, 2, 3, 4, 6},
			ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:     x509.KeyUsageDigitalSignature,
		}
	}

	// Create parent Certificate template
	if parentTemplate == nil {
		parentTemplate = &x509.Certificate{
			SerialNumber: big.NewInt(2019),
			Subject: pkix.Name{
				Organization:  []string{"ESTEEM"},
				Country:       []string{"Indo"},
				Province:      []string{"DKI"},
				Locality:      []string{"Jaksel"},
				StreetAddress: []string{"Jl Guru Mughni"},
				PostalCode:    []string{"94016"},
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().AddDate(10, 0, 0),
			IsCA:                  true,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			BasicConstraintsValid: true,
		}
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, certTemplate, parentTemplate, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, errors.New("GenerateCA " + err.Error())
	}

	// Encode Certificate
	publicKey := new(bytes.Buffer)
	pem.Encode(publicKey, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	return &CACertificate{
		PrivateKey:  caPrivKeyPEM.Bytes(),
		CSR:         csrPEM.Bytes(),
		Certificate: publicKey.Bytes(),
	}, nil
}

// Parse PrivateKey []byte into *rsa.PrivateKey
func ParsePrivateKey(privateKeyByte []byte, rsaPrivateKeyPassword string) (*rsa.PrivateKey, error) {
	privPem, _ := pem.Decode(privateKeyByte)
	var privPemBytes []byte
	var err error

	if rsaPrivateKeyPassword != "" {
		privPemBytes, err = x509.DecryptPEMBlock(privPem, []byte(rsaPrivateKeyPassword))
	} else {
		privPemBytes = privPem.Bytes
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil {
			log.Println("func GetCertificate, Unable to parse RSA private key")
			return nil, err
		}
	}

	var privateKey *rsa.PrivateKey
	var ok bool
	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("Unable to parse RSA private key")
	}

	return privateKey, nil
}

// Parse PublicKey []byte into *rsa.PublicKey
func ParsePublicKey(publicKeyByte []byte) (*rsa.PublicKey, error) {

	pubBlock, _ := pem.Decode(publicKeyByte)
	if pubBlock == nil {
		return nil, errors.New("ParseCSR, failed to parse csr PEM")

	}

	var csr *x509.Certificate
	var err error

	csr, err = x509.ParseCertificate(pubBlock.Bytes)
	if err != nil {
		return nil, errors.New("ParseCSR, failed to parse csr")
	}

	var publicKey *rsa.PublicKey
	var ok bool
	publicKey, ok = csr.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("Unable to parse RSA private key")
	}

	return publicKey, nil
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		log.Println(err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte, optionalPassword string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, []byte(optionalPassword))
		if err != nil {
			log.Println("BytesToPrivateKey ", err)
			return nil, err
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		log.Println("BytesToPrivateKey ", err)
		return nil, err
	}

	return key, nil
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Println("BytesToPublicKey ", err)
			return nil, err
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		log.Println(err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		log.Println("BytesToPrivateKey error parsing PublicKey")
		return nil, errors.New("error parsing PublicKey")
	}
	return key, nil
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		log.Println("EncryptWithPublicKey ", err)
		return nil, err
	}
	return ciphertext, nil
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		log.Println("DecryptWithPrivateKey ", err)
		return nil, err
	}
	return plaintext, nil
}
