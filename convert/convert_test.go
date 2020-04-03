package convert_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	. "github.com/zainul/ark/convert"
)

func TestToByteArr(t *testing.T) {
	assert.Equal(t, []byte(""), ToByteArr(nil))
	assert.Equal(t, []byte("A"), ToByteArr("A"))
	assert.Equal(t, []byte("3"), ToByteArr(3))
	assert.Equal(t, []byte("3"), ToByteArr(int64(3)))
	assert.Equal(t, []byte("false"), ToByteArr(false))
	assert.Equal(t, []byte("3.3E+00"), ToByteArr(3.3))
	assert.Equal(t, []byte("6"), ToByteArr([]byte("6")))
	assert.Equal(t, []byte("{}"), ToByteArr(errors.New("1")))

	jsonTest := struct {
		Test string
	}{}
	assert.Equal(t, []byte(`{"Test":""}`), ToByteArr(jsonTest))
	assert.Equal(t, []byte(""), ToByteArr(func() {}))

}

func TestToInt(t *testing.T) {
	assert.Equal(t, 6, ToInt("6"))
	assert.Equal(t, 6, ToInt(6))
	assert.Equal(t, 6, ToInt(int64(6)))
	assert.Equal(t, 6, ToInt(float64(6)))
	assert.Equal(t, 6, ToInt(uint8("6"[0])))
	assert.Equal(t, 6, ToInt([]byte("6")))
	assert.Equal(t, 6, ToInt([]interface{}{6}))
	assert.Equal(t, 0, ToInt('x'))
}

func TestToInt64(t *testing.T) {
	assert.Equal(t, int64(6), ToInt64("6"))
	assert.Equal(t, int64(6), ToInt64(6))
	assert.Equal(t, int64(6), ToInt64(int64(6)))
	assert.Equal(t, int64(6), ToInt64(float64(6)))
	assert.Equal(t, int64(6), ToInt64(uint8("6"[0])))
	assert.Equal(t, int64(6), ToInt64([]byte("6")))
	assert.Equal(t, int64(0), ToInt64('x'))

}

func TestToString(t *testing.T) {
	assert.Equal(t, "", ToString(nil))
	assert.Equal(t, "100", ToString("100"))
	assert.Equal(t, "100", ToString(100))
	assert.Equal(t, "100", ToString(int64(100)))
	assert.Equal(t, "6", ToString(uint8("6"[0])))
	assert.Equal(t, "false", ToString(false))
	assert.Equal(t, "100", ToString(float64(100)))
	assert.Equal(t, "100", ToString([]byte("100")))

	jsonTest := struct {
		Test string
	}{}
	assert.Equal(t, `{"Test":""}`, ToString(jsonTest))
	assert.Equal(t, "", ToString(func() {}))
}

func TestFixMySQLTime(t *testing.T) {
	const jakartaTZ = "Asia/Jakarta"
	var jakartaLocation, _ = time.LoadLocation(jakartaTZ)

	now := time.Now()

	expected := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond(),
		jakartaLocation)

	// mysql does not store timezone, so we assume the timezone is UTC
	tt := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond(),
		time.UTC)

	tcs := []struct {
		name   string
		input  time.Time
		offset int
	}{
		{
			name:   "utc-7",
			input:  tt.Add(time.Duration(-14) * time.Hour),
			offset: -7,
		},
		{
			name:   "utc",
			input:  tt.Add(time.Duration(-7) * time.Hour),
			offset: 0,
		},
		{
			name:   "utc+7",
			input:  tt,
			offset: 7,
		},
		{
			name:   "utc+8",
			input:  tt.Add(time.Duration(1) * time.Hour),
			offset: 8,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res := FixMySQLTime(tc.input, tc.offset*3600)
			assert.Equal(t, expected, res)
		})
	}
}

func TestToFloat64(t *testing.T) {
	assert.Equal(t, float64(6), ToFloat64("6"))
	assert.Equal(t, float64(6), ToFloat64(6))
	assert.Equal(t, float64(6), ToFloat64(int64(6)))
	assert.Equal(t, float64(6), ToFloat64(float64(6)))
	assert.Equal(t, float64(6), ToFloat64(uint8("6"[0])))
	assert.Equal(t, float64(0), ToFloat64([]byte("6")))

}
