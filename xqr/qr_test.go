package xqr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateQR(t *testing.T) {
	code := NewQR("hi toped", 300)

	err := code.GenerateQrCodeImage("toped.png")

	assert.Equal(t, nil, err)

	bt, err := code.GenerateQrCodeImageByte()

	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, bt)
}
