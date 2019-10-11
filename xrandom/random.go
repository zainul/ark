package xrandom

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"strconv"
	"strings"
	"time"

	mathRand "math/rand"
)

const (
	pkg       = "pkg.random"
	maxBigInt = 9223372036854775807
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	id := strings.ToUpper(base64.URLEncoding.EncodeToString(b))
	id = strings.Replace(id, "_", "", -1)
	id = strings.Replace(id, "-", "", -1)
	id = strings.Replace(id, "=", "", -1)
	return id, err
}

// GenerateRandomStringClean ...
func GenerateRandomStringClean(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	bs := base64.URLEncoding.EncodeToString(b)
	return strings.ToUpper(bs[:30]), err
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// GenerateAlias ...
func GenerateAlias(partWord string) string {
	if len(partWord) >= 10 {
		return partWord[:10]
	}

	partWord = partWord + Reverse(strconv.Itoa(int(time.Now().UnixNano()))[10-len(partWord):])

	return strings.ToUpper(partWord[:10])
}

// GenerateRandomInt ...
func GenerateRandomInt() (int64, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		return time.Now().Unix(), nil
	}

	n := nBig.Int64()

	return n, nil

}

// GenerateAliasPlusTime ...
func GenerateAliasPlusTime(alias string) string {
	strTime := time.Now().Format("060102150405") + strconv.Itoa(mathRand.Intn(1000000000))
	return strings.ToUpper((alias + strTime)[:24])
}

// GenerateGeneralID ...
func GenerateGeneralID() string {
	str, _ := GenerateRandomString(5)
	str1, _ := GenerateRandomString(5)
	str2, _ := GenerateRandomString(5)

	strTime := str + time.Now().Format("2006") + time.Now().Format("0102") + str1 + time.Now().Format("150405") + str2

	id := strings.Replace(strTime, "_", "", -1)
	id = strings.Replace(id, "-", "", -1)
	id = strings.Replace(id, "=", "", -1)

	str4, _ := GenerateRandomString(35)
	idx := id + str4

	return idx[:50]
}
