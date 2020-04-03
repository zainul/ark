package dummyelastic

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/zainul/ark/storage/elastic"
	eslib "gopkg.in/olivere/elastic.v5"
)

// New dummy module
func New(config Config) elastic.Elastic {
	return &dummy{
		config: config,
	}
}

// Connect to elastic
func (e *dummy) Connect() error {
	return nil
}

// GetClient of the elastic
func (e *dummy) GetClient() interface{} {
	return nil
}

// Search function
func (e *dummy) Search(ctx context.Context, q elastic.SearchParam) ([]string, []json.RawMessage, error) {
	// Validate empty
	if q.Request == nil {
		return nil, nil, errors.New("Invalid Search parameter")
	}

	request, _ := q.Request.(string)
	resultJSON, ok := e.config.Mocking[request]

	// No result
	if !ok {
		return nil, nil, nil
	}

	indexes := []string{}
	results := []json.RawMessage{}
	json.Unmarshal([]byte(resultJSON), &results)
	for i := range results {
		indexes = append(indexes, (q.Index + "_" + strconv.Itoa(i)))
	}
	return indexes, results, nil

}

func (e *dummy) Delete(ctx context.Context, param elastic.DeleteParam) error {
	if param.Index == "" || param.Type == "" || param.ID == "" {
		return errors.New("Invalid Delete Param")
	}

	return nil
}

// Suggest function
func (e *dummy) Suggest(ctx context.Context, q elastic.SuggestParam) ([]string, []json.RawMessage, error) {

	// Validate empty
	if q.Request == nil {
		return nil, nil, errors.New("Invalid Suggest parameter")
	}
	request := q.Request.(string)
	resultJSON, ok := e.config.Mocking[request]

	// No result
	if !ok {
		return nil, nil, nil
	}

	indexes := []string{}
	results := []json.RawMessage{}
	json.Unmarshal([]byte(resultJSON), &results)
	for i := range results {
		indexes = append(indexes, (q.Index + "_" + strconv.Itoa(i)))
	}
	return indexes, results, nil

}

// Get function
func (e *dummy) Get(ctx context.Context, q elastic.GetParam, target interface{}) error {
	// validate empty param
	if q.Index == "" || q.ID == "" {
		return errors.New("Invalid Suggest parameter")
	}

	resultJSON, ok := e.config.Mocking[q.ID]
	// No result
	if !ok {
		return nil
	}

	return json.Unmarshal([]byte(resultJSON), &target)
}

// BulkInsert function
func (e *dummy) BulkInsert(ctx context.Context, requestParams []elastic.InsertParam) error {
	for _, param := range requestParams {
		if param.ID == "testerror" {
			return errors.New("Foo Bar")
		}
	}
	return nil
}

// Upsert function
func (e *dummy) Upsert(ctx context.Context, requestParam elastic.InsertParam) (elastic.Response, error) {

	result := elastic.Response{}
	if requestParam.ID == "" {
		return result, errors.New("Foo Bar")
	}

	res, ok := e.config.Mocking[requestParam.ID]
	// Insert error, not registered in config.Mocking
	if !ok {
		return result, errors.New("Insert error")
	}

	err := json.Unmarshal([]byte(res), &result)

	return result, err
}

func (e *dummy) GetBucket(ctx context.Context, param elastic.GetBucketParam) (*eslib.AggregationBucketKeyItems, error) {
	var results *eslib.AggregationBucketKeyItems
	// Validate params
	if param.Index == "" || param.Source == "" || param.Terms == "" || param.Type == "" {
		return results, errors.New("Invalid parameter")
	}
	sourceQuery := param.Source.(string)
	resultJSON, ok := e.config.Mocking[sourceQuery]

	// No result
	if !ok {
		return results, nil
	}

	err := json.Unmarshal([]byte(resultJSON), &results)
	return results, err
}

func (e *dummy) Refresh(index string) {
}

func (e *dummy) Aggregation(ctx context.Context, param elastic.SearchParam) (map[string]eslib.AggregationBucketKeyItems, error) {
	const funcName = "Aggregation"

	if param.Index == "" || param.Request == nil {
		return nil, errors.New("Invalid Search parameter")
	}

	req, _ := param.Request.(string)
	if err, exists := e.config.ErrorMock[req]; exists {
		return nil, err
	}

	res, exists := e.config.Mocking[req]
	if !exists {
		return nil, errors.New("No mock available")
	}

	var searchRes eslib.SearchResult
	if err := json.Unmarshal([]byte(res), &searchRes); err != nil {
		return nil, err
	}

	aggregations := map[string]eslib.AggregationBucketKeyItems{}
	for aggKey, data := range searchRes.Aggregations {
		if data == nil {
			continue
		}

		var agg eslib.AggregationBucketKeyItems
		if err := json.Unmarshal(*data, &agg); err != nil {
			log.Printf("[%s] error unmarshal: %+v\n", funcName, err)
			continue
		}
		agg.Aggregations = nil
		for idx := range agg.Buckets {
			bucket := agg.Buckets[idx]
			if bucket == nil {
				continue
			}
			bucket.Aggregations = nil
		}
		aggregations[aggKey] = agg
	}

	return aggregations, nil
}
