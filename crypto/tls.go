package crypto

import (
	"crypto/tls"
	"crypto/x509"
	"log"
)

type TlsConfig struct {
	ServerConfig *tls.Config
	ClientConfig *tls.Config
}

func NewTlsConfig(privateKey []byte, publicKey []byte) *TlsConfig {
	serverCert, err := tls.X509KeyPair(publicKey, privateKey)
	if err != nil {
		log.Println("NewTlsConfig ", err)
		return nil
	}

	serverTLSConf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}

	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(publicKey)
	clientTLSConf := &tls.Config{
		RootCAs: certpool,
	}

	return &TlsConfig{
		ServerConfig: serverTLSConf,
		ClientConfig: clientTLSConf,
	}
}
