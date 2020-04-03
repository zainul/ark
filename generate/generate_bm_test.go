package generate_test

import (
	"testing"

	. "github.com/zainul/ark/generate"
)

func BenchmarkUUID(b *testing.B) {

	for n := 0; n < b.N; n++ {
		UUID()
	}

}

func BenchmarkMD5(b *testing.B) {

	for n := 0; n < b.N; n++ {
		MD5("a")
	}

}
func BenchmarkSHA1(b *testing.B) {

	for n := 0; n < b.N; n++ {
		SHA1("message", "secret")
	}

}

func BenchmarkRandomString10Alphamumeric(b *testing.B) {

	for n := 0; n < b.N; n++ {
		RandomString(10, StringAlphaNumeric)
	}

}

func BenchmarkRandomString100Alphamumeric(b *testing.B) {

	for n := 0; n < b.N; n++ {
		RandomString(100, StringAlphaNumeric)
	}

}

func BenchmarkRandomString100Alpha(b *testing.B) {

	for n := 0; n < b.N; n++ {
		RandomString(100, "")
	}

}
