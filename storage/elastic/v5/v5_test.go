package v5_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zainul/ark/storage/elastic"
	. "github.com/zainul/ark/storage/elastic/v5"
	eslib "gopkg.in/olivere/elastic.v5"
)

func TestConnectError(t *testing.T) {

	config := Config{}
	es := New(config)

	assert.EqualError(t, es.Connect(), "no active connection found: no Elasticsearch node available")
	assert.Nil(t, es.GetClient())
}

func TestConnectSuccess(t *testing.T) {

	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `{
			"name" : "elastic",
			"cluster_name" : "elastic",
			"cluster_uuid" : "XqzEn1WWQMa6dGZmdS3A1g",
			"version" : {
			  "number" : "5.2.2",
			  "build_hash" : "f9d9b74",
			  "build_date" : "2017-02-24T17:26:45.835Z",
			  "build_snapshot" : false,
			  "lucene_version" : "6.4.1"
			},
			"tagline" : "You Know, for Search"
		  }
		  `
		w.Write([]byte(resp))
	}

	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)

	assert.Nil(t, es.Connect())

	assert.NotNil(t, es.GetClient())
}

func TestSearchQueryError(t *testing.T) {

	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `{"name":"elastic","cluster_name":"elastic","cluster_uuid":"XqzEn1WWQMa6dGZmdS3A1g","version":{"number":"5.2.2","build_hash":"f9d9b74","build_date":"2017-02-24T17:26:45.835Z","build_snapshot":false,"lucene_version":"6.4.1"},"tagline":"You Know, for Search"}`

		if r.URL.String() == "/halo/halo/_search" {
			resp = `boom`
		}

		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()

	// Invalid param
	param := elastic.SearchParam{}
	_, _, err := es.Search(context.TODO(), param)
	assert.EqualError(t, err, "Invalid Search parameter")

	// Error during hit elastic
	param = elastic.SearchParam{
		Index:   "halo",
		Request: "foo",
	}
	index, result, err := es.Search(context.TODO(), param)
	assert.Empty(t, index)
	assert.Empty(t, result)
	assert.EqualError(t, err, "invalid character 'b' looking for beginning of value")
}

func TestSearchQueryEmptySuccess(t *testing.T) {

	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {

		resp := `{"name":"elastic","cluster_name":"elastic","cluster_uuid":"XqzEn1WWQMa6dGZmdS3A1g","version":{"number":"5.2.2","build_hash":"f9d9b74","build_date":"2017-02-24T17:26:45.835Z","build_snapshot":false,"lucene_version":"6.4.1"},"tagline":"You Know, for Search"}`

		if r.URL.String() == "/halo/halo/_search" {
			resp = `{
				"took": 4,
				"timed_out": false,
				"_shards": {
					"total": 5,
					"successful": 5,
					"failed": 0
				},
				"hits": {
					"total": 0,
					"max_score": null,
					"hits": []
				}
			}`
		}

		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()

	// Error during hit elastic
	param := elastic.SearchParam{
		Index:   "halo",
		Request: "foo",
	}
	index, result, err := es.Search(context.TODO(), param)
	assert.Nil(t, err)
	assert.Empty(t, result)
	assert.Empty(t, index)
}

func TestSearchQuerySuccess(t *testing.T) {

	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {

		resp := `{"name":"elastic","cluster_name":"elastic","cluster_uuid":"XqzEn1WWQMa6dGZmdS3A1g","version":{"number":"5.2.2","build_hash":"f9d9b74","build_date":"2017-02-24T17:26:45.835Z","build_snapshot":false,"lucene_version":"6.4.1"},"tagline":"You Know, for Search"}`

		if r.URL.String() == "/halo/halo/_search" {
			resp = `{
				"took": 1,
				"timed_out": false,
				"_shards": {
					"total": 5,"successful": 5,"skipped": 0,"failed": 0
				},
				"hits": {
					"total": 2,"max_score": 1,
					"hits": [
						{
							"_index": "halo",
							"_type": "halo",
							"_id": "35b4f0710c0d0e1b3fda664a27f27f38",
							"_score": 1,
							"_source": {
								"key": "key0",
								"agent_id": 0					
							}
						},
						{
							"_index": "halo",
							"_type": "halo",
							"_id": "c8e0f4477f8c2884010e697dd460a70c",
							"_score": 1,
							"_source": {
								"key": "key1",
								"agent_id": 1
							}
						}
					]
				}
			}`
		}

		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()
	indexes := []string{"35b4f0710c0d0e1b3fda664a27f27f38", "c8e0f4477f8c2884010e697dd460a70c"}
	// Error during hit elastic
	param := elastic.SearchParam{
		Index:   "halo",
		Request: "foo",
	}
	index, result, err := es.Search(context.TODO(), param)
	assert.Nil(t, err)
	assert.NotEmpty(t, result)
	for k, v := range result {
		res := dummyResult{}
		json.Unmarshal(v, &res)

		assert.Equal(t, k, res.AgentID)
		assert.Equal(t, index[k], indexes[k])
		assert.Equal(t, fmt.Sprintf("key%d", k), res.Key)

	}
}

func TestSuggestQueryError(t *testing.T) {

	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `{"name":"elastic","cluster_name":"elastic","cluster_uuid":"XqzEn1WWQMa6dGZmdS3A1g","version":{"number":"5.2.2","build_hash":"f9d9b74","build_date":"2017-02-24T17:26:45.835Z","build_snapshot":false,"lucene_version":"6.4.1"},"tagline":"You Know, for Search"}`

		if r.URL.String() == "/halo/halo/_search" {
			resp = `boom`
		}

		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()

	// Invalid param
	param := elastic.SuggestParam{}
	_, _, err := es.Suggest(context.TODO(), param)
	assert.EqualError(t, err, "Invalid Suggest parameter")

	// Error during hit elastic
	param = elastic.SuggestParam{
		Index:       "halo",
		Request:     "foo",
		SuggestName: "bar",
	}
	index, result, err := es.Suggest(context.TODO(), param)
	assert.Empty(t, index)
	assert.Empty(t, result)
	assert.EqualError(t, err, "invalid character 'b' looking for beginning of value")
}

func TestSuggestQueryEmptySuccess(t *testing.T) {

	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {

		resp := `{"name":"elastic","cluster_name":"elastic","cluster_uuid":"XqzEn1WWQMa6dGZmdS3A1g","version":{"number":"5.2.2","build_hash":"f9d9b74","build_date":"2017-02-24T17:26:45.835Z","build_snapshot":false,"lucene_version":"6.4.1"},"tagline":"You Know, for Search"}`

		if r.URL.String() == "/halo/halo/_search" {
			resp = `{
				"took": 4,
				"timed_out": false,
				"_shards": {
					"total": 5,
					"successful": 5,
					"failed": 0
				},
				"hits": {
					"total": 0,
					"max_score": 0,
					"hits": []
				},
				"suggest": {
					"search-suggest": [
						{
							"text": "xxxxxxxxxxxxxxxx",
							"offset": 0,
							"length": 9,
							"options": []
						}
					]
				}
			}`
		}

		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()

	// Error during hit elastic
	param := elastic.SuggestParam{
		SuggestName: "search-suggest",
		Index:       "halo",
		Request:     "foo",
	}
	index, result, err := es.Suggest(context.TODO(), param)
	assert.Nil(t, err)
	assert.Empty(t, index)
	assert.Empty(t, result)
}

func TestSuggestQuerySuccess(t *testing.T) {

	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {

		resp := `{"name":"elastic","cluster_name":"elastic","cluster_uuid":"XqzEn1WWQMa6dGZmdS3A1g","version":{"number":"5.2.2","build_hash":"f9d9b74","build_date":"2017-02-24T17:26:45.835Z","build_snapshot":false,"lucene_version":"6.4.1"},"tagline":"You Know, for Search"}`

		if r.URL.String() == "/halo/halo/_search" {
			resp = `{
				"took": 150,
				"timed_out": false,
				"_shards": {
					"total": 5,
					"successful": 5,
					"failed": 0
				},
				"hits": {
					"total": 0,
					"max_score": 0,
					"hits": []
				},
				"suggest": {
					"search-suggest": [
						{
							"text": "agent",
							"offset": 0,
							"length": 3,
							"options": [
								{
									"text": "Agent1",
									"_index": "halo",
									"_type": "halo",
									"_id": "2",
									"_score": 152,
									"_source": {
										"key": "key0",
										"agent_id": 0
									}
								}
							]
						},
						{
							"text": "agent",
							"offset": 0,
							"length": 3,
							"options": [
								{
									"text": "Agent1",
									"_index": "halo",
									"_type": "halo",
									"_id": "3",
									"_score": 152,
									"_source": {
										"key": "key1",
										"agent_id": 1
									}
								}
							]
						}
					]
				}
			}`
		}

		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()
	ids := []string{"2", "3"}
	// Error during hit elastic
	param := elastic.SuggestParam{
		SuggestName: "search-suggest",
		Index:       "halo",
		Request:     "foo",
	}
	index, result, err := es.Suggest(context.TODO(), param)
	assert.Nil(t, err)
	assert.NotEmpty(t, result)
	for k, v := range result {
		res := dummyResult{}
		json.Unmarshal(v, &res)

		assert.Equal(t, k, res.AgentID)
		assert.Equal(t, fmt.Sprintf("key%d", k), res.Key)
		assert.Equal(t, ids[k], index[k])

	}
}

func TestGetQuerySuccess(t *testing.T) {
	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `{"_index":"halo","_type":"halo","_id":"halohalo","_version":21,"found":true,"_source":{"key":"haloe","agent_id":3}}`
		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()

	// Error during hit elastic
	param := elastic.GetParam{
		Index:  "halo",
		ID:     "halohalo",
		Source: []string{"key", "agent_id"},
	}

	responseResult := dummyResult{}
	err := es.Get(context.TODO(), param, &responseResult)
	assert.Nil(t, err)
	assert.NotEmpty(t, responseResult)
}

func TestGetQueryNotFound(t *testing.T) {
	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `{"_index":"halo","_type":"halo","_id":"halohalo","_version":21,"found":false,"_source":{}}`
		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()

	// Error during hit elastic
	param := elastic.GetParam{
		Index: "halo",
		ID:    "halohalo",
	}

	responseResult := dummyResult{}
	err := es.Get(context.TODO(), param, &responseResult)
	assert.Error(t, err)
	assert.Empty(t, responseResult)

	// Test not set Index
	param.Index = ""

	err = es.Get(context.TODO(), param, &responseResult)
	assert.EqualError(t, err, IndexIsEmpty)

	// Test not set ID
	param.Index = "halo"
	param.ID = ""

	err = es.Get(context.TODO(), param, &responseResult)
	assert.Error(t, err)

}

func TestGetFailedParse(t *testing.T) {
	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `{"error":{"root_cause":[{"type":"illegal_argument_exception","reason":"No endpoint or operation is available at [travel_airlines_availability]"}],"type":"illegal_argument_exception","reason":"No endpoint or operation is available at [travel_airlines_availability]"},"status":400}`
		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()

	// Error during hit elastic
	param := elastic.GetParam{
		Index: "halo",
		ID:    "halohalo",
	}

	responseResult := dummyResult{}
	err := es.Get(context.TODO(), param, " ")

	assert.Error(t, err)
	assert.Empty(t, responseResult)

	// Test result is nil
	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := ``
		w.Write([]byte(resp))
	}

	err = es.Get(context.TODO(), param, &responseResult)

	assert.Error(t, err)
	assert.Empty(t, responseResult)
}

func TestBulkInsertSuccess(t *testing.T) {
	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `
			{
				"took": 30,
				"errors": false,
				"items": [
					{
						"index": {
							"_index": "Test",
							"_type": "Test",
							"_id": "1",
							"_version": 1,
							"result": "created",
							"_shards": {
								"total": 2,
								"successfull": 1,
								"failed": 0
							},
							"status": 201,
							"_seq_no": 0,
							"_primary_term": 1
						}
					}
				]
			}
		`
		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()

	param := elastic.InsertParam{
		Index: "Test",
		ID:    "1",
		Doc:   "test",
	}
	params := make([]elastic.InsertParam, 0)
	params = append(params, param)

	err := es.BulkInsert(context.TODO(), params)

	assert.Nil(t, err)
}

func TestBulkInsertError(t *testing.T) {
	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `
			
		`
		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()

	param := elastic.InsertParam{
		Index: "Test",
		ID:    "1",
		Doc:   "test",
	}
	params := make([]elastic.InsertParam, 0)
	params = append(params, param)

	err := es.BulkInsert(context.TODO(), params)

	assert.NotNil(t, err)
}

func TestInsertError(t *testing.T) {
	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `{"_index":"travel_feeds","_type":"doc","_id":"2572","_version":5,"result":"updated","_shards":{"total":2,"successful":1,"failed":0}}`
		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()

	param := elastic.InsertParam{
		ID: "",
	}
	_, err := es.Upsert(context.TODO(), param)

	assert.NotNil(t, err)

}

func TestInsertSuccess(t *testing.T) {
	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `{"_index":"travel_feeds","_type":"doc","_id":"2572","_version":5,"result":"updated","_shards":{"total":2,"successful":1,"failed":0}}`
		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()

	param := elastic.InsertParam{
		ID:    "1",
		Type:  "doc",
		Index: "test",
		Doc:   "{}",
	}
	_, err := es.Upsert(context.TODO(), param)

	assert.Nil(t, err)
}

func TestRefresh(t *testing.T) {

	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `{"_index":"travel_feeds","_type":"doc","_id":"2572","_version":5,"result":"updated","_shards":{"total":2,"successful":1,"failed":0}}`
		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()
	es.Refresh("travel_feeds")
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
		mockServer     func() *httptest.Server
	}

	testCases := []args{
		args{
			name:          "Invalid params",
			expectedError: errors.New("Invalid Search parameter"),
		},
		args{
			name: "Error response",
			param: elastic.SearchParam{
				Index:   "property",
				Request: `{}`,
			},
			expectedError: &eslib.Error{Status: 404},
			mockServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
					http.NotFound(writer, req)
				}))
			},
		},
		args{
			name: "Success",
			param: elastic.SearchParam{
				Index: "travel_property",
				Request: `
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
				`,
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
			mockServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
					func(writer http.ResponseWriter, req *http.Request) {
						res := `
						{
							"took": 3,
							"timed_out": false,
							"_shards": {
							  "total": 5,
							  "successful": 5,
							  "skipped": 0,
							  "failed": 0
							},
							"hits": {
							  "total": 3319,
							  "max_score": 0.0,
							  "hits": []
							},
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
						`

						writer.Write([]byte(res))
					}(writer, req)
				}))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.mockServer != nil {
				httpServer := testCase.mockServer()
				testCase.config.Endpoint = httpServer.URL
				defer httpServer.Close()
			}

			es := New(testCase.config)
			aggs, err := es.Aggregation(context.Background(), testCase.param)
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult, aggs)
		})
	}
}
