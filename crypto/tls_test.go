package crypto_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zainul/ark/crypto"
)

func TestNewTlsConfig(t *testing.T) {
	ca, err := crypto.GenerateCA(nil, nil, nil)

	// get our ca and server certificate
	tlsConfig := crypto.NewTlsConfig(ca.PrivateKey, ca.Certificate)
	assert.NotNil(t, tlsConfig)
	if tlsConfig == nil {
		return
	}

	// set up the httptest.Server using our certificate signed by our CA
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "success!")
	}))

	server.TLS = tlsConfig.ServerConfig
	server.StartTLS()
	defer server.Close()

	// communicate with the server using an http.Client configured to trust our CA
	transport := &http.Transport{
		TLSClientConfig: tlsConfig.ClientConfig,
	}
	http := http.Client{
		Transport: transport,
	}
	resp, err := http.Get(server.URL)
	assert.Nil(t, err)

	// verify the response
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	body := strings.TrimSpace(string(respBodyBytes[:]))
	assert.Equal(t, "success!", body)
}
