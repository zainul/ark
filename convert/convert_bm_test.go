package convert_test

import (
	"testing"

	. "github.com/zainul/ark/convert"
)

func BenchmarkToFloat64(b *testing.B) {

	for n := 0; n < b.N; n++ {
		ToFloat64(2)
	}

}

func BenchmarkToByteArr(b *testing.B) {

	for n := 0; n < b.N; n++ {
		ToByteArr("4545")
	}

}
func BenchmarkToInt(b *testing.B) {

	for n := 0; n < b.N; n++ {
		ToInt("4545")
	}

}
