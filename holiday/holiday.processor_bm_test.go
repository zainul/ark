package holiday

import (
	"testing"
)

func BenchmarkGet(b *testing.B) {

	data := &Config{
		Holidays: dummyHolidays,
	}

	p := GetProcessor(data)

	for n := 0; n < b.N; n++ {
		p.Get("ID")
	}

}
