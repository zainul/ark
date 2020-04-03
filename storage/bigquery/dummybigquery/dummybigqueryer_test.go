package dummybigquery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDummy(t *testing.T) {

	// mock error connection here
	newArgs := DummyArgs{
		ExpectedResult:    []string{`{"count":0,"user":1}`},
		IsConnectionError: true,
		CurrentIndex:      0,
	}

	newWrongBQ := New(newArgs)
	err := newWrongBQ.Connect()
	assert.Error(t, err)

	// mock success connection here
	newArgs = DummyArgs{
		ExpectedResult:    []string{`{"count":0,"user":1}`},
		IsConnectionError: false,
		CurrentIndex:      0,
	}

	newDummyBQ := New(newArgs)
	err = newDummyBQ.Connect()
	assert.Nil(t, err)

	queryResult := newDummyBQ.Query("SELECT * FROM ")

	var targetTotal []Target
	for {
		var target Target
		err := queryResult.Next(&target)
		if err != nil {
			break
		}

		targetTotal = append(targetTotal, target)
	}
	assert.NotNil(t, targetTotal)

}
