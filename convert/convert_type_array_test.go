package convert_test

import (
	"encoding/json"
	"github.com/zainul/ark/convert"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	A    string `json:"a"`
	Subs []Sub  `json:"subs"`
}

type ComparedTestStruct struct {
	ACompare    string `json:"a"`
	SubsCompare []Sub  `json:"subs"`
}

type Sub struct {
	Sub int `json:"sub"`
}

func TestIntArrToStringArr(t *testing.T) {
	a := []int{1, 2, 3, 4}
	assert.Equal(t, []string{"1", "2", "3", "4"}, convert.IntArrToStringArr(a))
}

func TestParseToDestination(t *testing.T) {
	obj := ComparedTestStruct{
		ACompare:    "huruf A",
		SubsCompare: []Sub{
			{
				1,
			},
			{
				2,
			},
		},
	}

	expected := TestStruct{
		A: "huruf A",
		Subs: []Sub{
			{
				1,
			},
			{
				2,
			},
		},
	}

	var result TestStruct

	err := convert.ParseToDestination(obj, &result)
	assert.Equal(t, expected, result)
	assert.Equal(t, nil, err)

	err = convert.ParseToDestination(`{`, `{`)
	_, ok := err.(*json.InvalidUnmarshalError)
	assert.True(t, ok)
	err = convert.ParseToDestination(math.Inf(1), ``)
	_, ok = err.(*json.UnsupportedValueError)
	assert.True(t, ok)
}
