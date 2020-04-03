package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zainul/ark/crypto"
)

func TestGenerateKeyPairByte(t *testing.T) {
	ca, err := crypto.GenerateCA(nil, nil, nil)
	assert.Nil(t, err)

	if err != nil {
		return
	}

	priv, err := crypto.ParsePrivateKey(ca.PrivateKey, "")
	assert.Nil(t, err)
	pub, err := crypto.ParsePublicKey(ca.Certificate)
	assert.Nil(t, err)

	if err == nil {
		plainText := "HALOOO DUNIA"
		cipherText, _ := crypto.EncryptWithPublicKey([]byte(plainText), pub)

		decrypted, _ := crypto.DecryptWithPrivateKey(cipherText, priv)

		// Test if decrypted file match the original plainText
		assert.Equal(t, plainText, string(decrypted))
	}
}
func TestEncryptDecrypt(t *testing.T) {
	var err error

	plainText := "HALOOO DUNIA"

	// Get Certificate
	ca, err := crypto.GenerateCA(nil, nil, nil)
	assert.Nil(t, err)

	if err != nil {
		return
	}

	pub, err := crypto.ParsePublicKey(ca.Certificate)
	priv, err := crypto.ParsePrivateKey(ca.PrivateKey, "")

	if err == nil {
		cipherText, _ := crypto.EncryptWithPublicKey([]byte(plainText), pub)
		assert.Nil(t, err)

		decrypted, _ := crypto.DecryptWithPrivateKey(cipherText, priv)

		// Test if decrypted file match the original plainText
		assert.Equal(t, plainText, string(decrypted))
	}
}
