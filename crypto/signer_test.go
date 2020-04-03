package crypto_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zainul/ark/crypto"
)

func TestSignWithPrivateKey(t *testing.T) {
	caCert, err := crypto.GenerateCA(nil, nil, nil)
	assert.Nil(t, err)

	if err != nil {
		return
	}

	privKey, err := crypto.ParsePrivateKey(caCert.PrivateKey, "")
	assert.Nil(t, err)
	pubKey, err := crypto.ParsePublicKey(caCert.Certificate)
	assert.Nil(t, err)

	if err == nil {

		signer := crypto.NewSignerFromKey(privKey)

		toSign := "date: Thu, 05 Jan 2012 21:31:40 GMT"

		signed, err := signer.SignWithSHA256([]byte(toSign))
		assert.Nil(t, err)

		sig := base64.StdEncoding.EncodeToString(signed)
		fmt.Printf("Signature: %v\n", sig)

		parser := crypto.NewUnsignerFromKey(pubKey)

		err = parser.UnsignWithSHA256([]byte(toSign), signed)
		assert.Nil(t, err)
	}
}
