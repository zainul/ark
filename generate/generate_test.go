package generate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/zainul/ark/generate"
)

func TestUUID(t *testing.T) {
	assert.Len(t, UUID(), 32)
}

func TestMD5(t *testing.T) {
	assert.Equal(t, "0cc175b9c0f1b6a831c399e269772661", MD5("a"))
}

func TestSHA1(t *testing.T) {
	assert.Equal(t, "DK9kn+7klT2Hv5A6wRdsReAo3xY=", SHA1("message", "secret"))
}

func TestRandomString(t *testing.T) {
	assert.Len(t, RandomString(10, StringAlpha), 10)
	assert.Len(t, RandomString(10, ""), 10)
}
