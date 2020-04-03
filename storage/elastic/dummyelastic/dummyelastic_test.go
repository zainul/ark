package dummyelastic_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zainul/ark/storage/elastic"
	. "github.com/zainul/ark/storage/elastic/dummyelastic"
	eslib "gopkg.in/olivere/elastic.v5"
)

func TestGetClientSuccess(t *testing.T) {

	config := Config{}
	es := New(config)

	assert.Nil(t, es.Connect())
	assert.Nil(t, es.GetClient())
}

func TestSearchQueryError(t *testing.T) {

	config := Config{}
	es := New(config)

	// Invalid param
	param := elastic.SearchParam{}
	index, result, err := es.Search(context.TODO(), param)
	assert.Empty(t, index)
	assert.Empty(t, result)
	assert.EqualError(t, err, "Invalid Search parameter")

}

func TestSearchQueryEmptySuccess(t *testing.T) {

	config := Config{}
	es := New(config)

	// Without mocking
	param := elastic.SearchParam{
		Request: "empty response",
	}
	index, result, err := es.Search(context.TODO(), param)
	assert.Empty(t, result)
	assert.Empty(t, index)
	assert.Nil(t, err)

}

func TestSearchQuerySuccess(t *testing.T) {

	mocking := map[string]string{}

	mocking["foo"] = `
	[{
		"key": "key0",
		"agent_id": 0					
	},
	{
		"key": "key1",
		"agent_id": 1
	}]`
	config := Config{
		Mocking: mocking,
	}
	es := New(config)

	// With mocking
	param := elastic.SearchParam{
		Index:   "bar",
		Request: "foo",
	}
	index, result, err := es.Search(context.TODO(), param)

	assert.Nil(t, err)
	for k, v := range result {
		res := dummyResult{}
		json.Unmarshal(v, &res)

		assert.Equal(t, k, res.AgentID)
		assert.Equal(t, index[k], param.Index+"_"+strconv.Itoa(k))
		assert.Equal(t, fmt.Sprintf("key%d", k), res.Key)

	}

}

func TestSuggestQueryError(t *testing.T) {

	config := Config{}
	es := New(config)

	// Invalid param
	param := elastic.SuggestParam{}
	index, result, err := es.Suggest(context.TODO(), param)
	assert.Empty(t, result)
	assert.Empty(t, index)
	assert.EqualError(t, err, "Invalid Suggest parameter")

}

func TestSuggestQueryEmptySuccess(t *testing.T) {

	config := Config{}
	es := New(config)

	// Without mocking
	param := elastic.SuggestParam{
		Request: "empty response",
	}
	index, result, err := es.Suggest(context.TODO(), param)
	assert.Empty(t, result)
	assert.Empty(t, index)
	assert.Nil(t, err)

}

func TestSuggestQuerySuccess(t *testing.T) {

	mocking := map[string]string{}

	mocking["foo"] = `
	[{
		"key": "key0",
		"agent_id": 0					
	},
	{
		"key": "key1",
		"agent_id": 1
	}]`
	config := Config{
		Mocking: mocking,
	}
	es := New(config)

	// With mocking
	param := elastic.SuggestParam{
		Request: "foo",
	}
	index, result, err := es.Suggest(context.TODO(), param)

	assert.Nil(t, err)
	for k, v := range result {
		res := dummyResult{}
		json.Unmarshal(v, &res)

		assert.Equal(t, k, res.AgentID)
		assert.Equal(t, fmt.Sprintf("key%d", k), res.Key)
		assert.Equal(t, index[k], param.Index+"_"+strconv.Itoa(k))

	}

}

func TestGetQuerySuccess(t *testing.T) {
	mocking := map[string]string{}

	mocking["halo"] = `
	{
		"key": "key1",
		"agent_id": 1
	}`
	config := Config{
		Mocking: mocking,
	}
	es := New(config)

	// With mocking
	param := elastic.GetParam{
		ID:    "halo",
		Index: "halo",
	}

	res := dummyResult{}
	err := es.Get(context.TODO(), param, &res)

	assert.Nil(t, err)
	assert.NotEmpty(t, res)
}

func TestGetQueryFailed(t *testing.T) {
	mocking := map[string]string{}

	mocking["asdq"] = `
	{
		"key": "key1",
		"agent_id": 1
	}`
	config := Config{
		Mocking: mocking,
	}
	es := New(config)

	// With mocking
	param := elastic.GetParam{
		ID:    "halo",
		Index: "halo",
	}

	res := dummyResult{}
	err := es.Get(context.TODO(), param, &res)

	assert.Nil(t, err)
	assert.Empty(t, res)

	param.ID = ""
	err = es.Get(context.TODO(), param, &res)
	assert.Error(t, err)
}

func TestBulkInsertSuccess(t *testing.T) {

	mocking := map[string]string{}

	mocking["foo"] = ``
	config := Config{
		Mocking: mocking,
	}
	es := New(config)

	// With mocking
	param := elastic.InsertParam{
		ID: "testsuccess",
	}
	params := make([]elastic.InsertParam, 0)
	params = append(params, param)
	err := es.BulkInsert(context.TODO(), params)

	assert.Nil(t, err)
}

func TestBulkInsertError(t *testing.T) {

	mocking := map[string]string{}

	mocking["foo"] = ``
	config := Config{
		Mocking: mocking,
	}
	es := New(config)

	// With mocking
	param := elastic.InsertParam{
		ID: "testerror",
	}
	params := make([]elastic.InsertParam, 0)
	params = append(params, param)
	err := es.BulkInsert(context.TODO(), params)

	assert.NotNil(t, err)
}

func TestInsertError(t *testing.T) {

	mocking := map[string]string{}

	mocking["foo"] = ``
	config := Config{
		Mocking: mocking,
	}
	es := New(config)

	// With mocking
	param := elastic.InsertParam{
		ID: "",
	}
	res, err := es.Upsert(context.TODO(), param)

	assert.NotNil(t, err)
	assert.Equal(t, elastic.Response{}, res)

	// Insert error
	param.ID = "1"
	res, err = es.Upsert(context.TODO(), param)

	assert.NotNil(t, err)
	assert.Equal(t, elastic.Response{}, res)
}

func TestInsertSuccess(t *testing.T) {

	mocking := map[string]string{}

	mocking["foo"] = `
	{
		"status":200
	}
	`
	config := Config{
		Mocking: mocking,
	}
	es := New(config)

	// With mocking
	param := elastic.InsertParam{
		ID: "foo",
	}
	res, err := es.Upsert(context.TODO(), param)

	assert.Nil(t, err)
	assert.Equal(t, 200, res.Status)
}

func TestRefresh(t *testing.T) {
	mocking := map[string]string{}
	config := Config{
		Mocking: mocking,
	}
	es := New(config)
	es.Refresh("test")
}

func TestGetBucket(t *testing.T) {
	type testCase struct {
		Name          string
		Mock          map[string]string
		Param         interface{}
		ExpectedError bool
	}
	testCases := []testCase{
		{
			Name: "Success",
			Mock: map[string]string{"test_source": `{}`},
			Param: elastic.GetBucketParam{
				Index:  "test_index",
				Type:   "test_type",
				Source: "test_source",
				Terms:  "test_terms",
			},
			ExpectedError: false,
		},
		{
			Name: "Mock not found",
			Mock: map[string]string{"xxx": ``},
			Param: elastic.GetBucketParam{
				Index:  "test_index",
				Type:   "test_type",
				Source: "test_source",
				Terms:  "test_terms",
			},
			ExpectedError: false,
		},
		{
			Name: "Invalid param",
			Mock: map[string]string{"test_source": `{}`},
			Param: elastic.GetBucketParam{
				Index:  "test_index",
				Type:   "test_type",
				Source: "",
				Terms:  "test_terms",
			},
			ExpectedError: true,
		},
		{
			Name: "Get bucket error",
			Mock: map[string]string{"test_source": ``},
			Param: elastic.GetBucketParam{
				Index:  "test_index",
				Type:   "test_type",
				Source: "test_source",
				Terms:  "test_terms",
			},
			ExpectedError: true,
		},
	}

	for _, testcase := range testCases {
		config := Config{
			Mocking: testcase.Mock,
		}
		es := New(config)
		_, err := es.GetBucket(context.TODO(), testcase.Param.(elastic.GetBucketParam))
		if testcase.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

type dummyResult struct {
	Key     string `json:"key"`
	AgentID int    `json:"agent_id"`
}

func TestAggregation(t *testing.T) {
	type args struct {
		name           string
		param          elastic.SearchParam
		expectedError  error
		expectedResult map[string]eslib.AggregationBucketKeyItems
		config         Config
	}

	query := `
	{
		"query": {
			"bool": {
				"must": [
					{"terms": {"location.city.id": [57]}}
				]
			}
		},
		"size": 0,
		"from": 0,
		"aggs": {
			"agg_chain_id": {
				"terms": {
					"field": "chain_id",
					"size": 5
				}
			}
		}
	}
	`
	testCases := []args{
		args{
			name:          "Invalid params",
			expectedError: errors.New("Invalid Search parameter"),
		},
		args{
			name: "Error response",
			param: elastic.SearchParam{
				Index:   "property",
				Request: query,
			},
			expectedError: errors.New("Error"),
			config: Config{
				ErrorMock: map[string]error{
					query: errors.New("Error"),
				},
			},
		},
		args{
			name: "Error no mock",
			param: elastic.SearchParam{
				Index:   "property",
				Request: query,
			},
			expectedError: errors.New("No mock available"),
			config:        Config{},
		},
		args{
			name: "Success",
			param: elastic.SearchParam{
				Index:   "travel_property",
				Request: query,
			},
			config: Config{
				Mocking: map[string]string{
					query: `
					{
						"aggregations": {
						  "xxx": null,
						  "yyy": "yyy",
						  "agg_chain_id": {
							"doc_count_error_upper_bound": 5,
							"sum_other_doc_count": 102,
							"buckets": [
							  null,
							  {
								"key": 0,
								"doc_count": 2828
							  },
							  {
								"key": 6015,
								"doc_count": 178
							  },
							  {
								"key": 6539,
								"doc_count": 153
							  },
							  {
								"key": 6627,
								"doc_count": 50
							  },
							  {
								"key": 1051,
								"doc_count": 8
							  }
							]
						  }
						}
					  }
					`,
				},
			},
			expectedResult: map[string]eslib.AggregationBucketKeyItems{
				"agg_chain_id": eslib.AggregationBucketKeyItems{
					DocCountErrorUpperBound: 5,
					SumOfOtherDocCount:      102,
					Buckets: []*eslib.AggregationBucketKeyItem{
						nil,
						&eslib.AggregationBucketKeyItem{
							Key:       float64(0),
							KeyNumber: json.Number("0"),
							DocCount:  2828,
						},
						&eslib.AggregationBucketKeyItem{
							Key:       float64(6015),
							KeyNumber: json.Number("6015"),
							DocCount:  178,
						},
						&eslib.AggregationBucketKeyItem{
							Key:       float64(6539),
							KeyNumber: json.Number("6539"),
							DocCount:  153,
						},
						&eslib.AggregationBucketKeyItem{
							Key:       float64(6627),
							KeyNumber: json.Number("6627"),
							DocCount:  50,
						},
						&eslib.AggregationBucketKeyItem{
							Key:       float64(1051),
							KeyNumber: json.Number("1051"),
							DocCount:  8,
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			es := New(testCase.config)
			aggs, err := es.Aggregation(context.Background(), testCase.param)
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult, aggs)
		})
	}
}
