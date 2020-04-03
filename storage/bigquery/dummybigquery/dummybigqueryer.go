package dummybigquery

import (
	"encoding/json"
	"errors"
	"github.com/zainul/ark/storage/bigquery"
	"reflect"
	"google.golang.org/api/iterator"
)

// New bigquery module
func New(mockParam DummyArgs) DummyBigQueryer {

	return &dummybigQueryModule{
		Parameters: mockParam,
	}
}

// Connect to oauth
func (b *dummybigQueryModule) Connect() error {

	if b.Parameters.IsConnectionError {
		return errors.New("error connection to BQ")
	}
	return nil
}

func (b *dummybigQueryModule) Query(queryParam string) bigquery.Iterator {
	return &IteratorDummy{
		ExpectedResult: b.Parameters.ExpectedResult,
		CurrentIndex:   b.Parameters.CurrentIndex,
	}

}

func (c *IteratorDummy) Next(dest interface{}) error {

	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}

	found := false
	var scannedVal string

	for m := range c.ExpectedResult {
		if m == c.CurrentIndex {
			found = true
			scannedVal = c.ExpectedResult[m]
			break
		}
	}

	if !found {
		return iterator.Done
		// return errors.New("cannot find expected index return no result exists")
	}

	var maps map[string]interface{}
	json.Unmarshal([]byte(scannedVal), &maps)
	s := v.Elem()

	for i := 0; i < s.NumField(); i++ {
		tf := s.Type().Field(i)
		tag := tf.Tag.Get("bigquery")


		if tag == "" {
			continue
		}

		val := maps[tag]
		if val == nil {
			continue
		}

		switch s.Field(i).Kind() {
		case reflect.Int, reflect.Int64:
			s.Field(i).SetInt(int64(val.(float64)))
		case reflect.String:
			s.Field(i).SetString(val.(string))
		case reflect.Float64:
			s.Field(i).SetFloat(val.(float64))
		default:
			break
		}
	}

	/*
		if err := json.Unmarshal([]byte(scannedVal), &dest); err != nil {
			return errors.New("cannot scan expected value into destination object")
		}
	*/
	c.CurrentIndex++

	return nil
}
