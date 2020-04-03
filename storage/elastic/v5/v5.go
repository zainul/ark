package v5

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/zainul/ark/storage/elastic"
	eslib "gopkg.in/olivere/elastic.v5"
)

// New v5 elastic module
func New(config Config) elastic.Elastic {
	return &elasticV5{
		config: config,
	}
}

// Connect to elastic
func (e *elasticV5) Connect() error {

	var err error

	// Connect to elastic client
	e.client, err = eslib.NewClient(
		eslib.SetURL(e.config.Endpoint),
		eslib.SetSniff(false),
		eslib.SetHealthcheck(false))

	return err
}

// GetClient of the elastic
func (e *elasticV5) GetClient() interface{} {
	return e.client
}

// Search function
func (e *elasticV5) Search(ctx context.Context, q elastic.SearchParam) ([]string, []json.RawMessage, error) {
	err := e.Connect()
	if err != nil {
		return nil, nil, err
	}

	// Validate empty
	if q.Index == "" || q.Request == nil {
		return nil, nil, errors.New("Invalid Search parameter")
	}

	if e.client == nil {
		return nil, nil, NilClientError
	}

	searchService := e.client.Search().Index(q.Index).Type(q.Index)
	queryResult, err := searchService.Source(q.Request).Do(ctx)
	if err != nil {
		return nil, nil, err
	}
	elasticIndexes := []string{}
	results := []json.RawMessage{}
	for _, h := range queryResult.Hits.Hits {
		results = append(results, *h.Source)
		elasticIndexes = append(elasticIndexes, h.Id)
	}

	return elasticIndexes, results, nil

}

// Suggest function
func (e *elasticV5) Suggest(ctx context.Context, q elastic.SuggestParam) ([]string, []json.RawMessage, error) {
	err := e.Connect()
	if err != nil {
		return nil, nil, err
	}

	// Validate empty
	if q.Index == "" || q.Request == nil || q.SuggestName == "" {
		return nil, nil, errors.New("Invalid Suggest parameter")
	}

	if e.client == nil {
		return nil, nil, NilClientError
	}

	searchService := e.client.Search().Index(q.Index).Type(q.Index)
	queryResult, err := searchService.Source(q.Request).Do(ctx)
	if err != nil {
		return nil, nil, err
	}

	results := []json.RawMessage{}
	elasticIndexes := []string{}

	// Process elastic suggestions
	if suggestions, ok := queryResult.Suggest[q.SuggestName]; ok {

		for r := range suggestions {

			for o := range suggestions[r].Options {
				results = append(results, *suggestions[r].Options[o].Source)
				elasticIndexes = append(elasticIndexes, suggestions[r].Options[o].Id)

			}

		}

	}

	return elasticIndexes, results, nil

}

// Get function
func (e *elasticV5) Get(ctx context.Context, q elastic.GetParam, target interface{}) error {

	// validate empty param
	if q.Index == "" {
		return errors.New(IndexIsEmpty)
	}

	if q.ID == "" {
		return errors.New(IDIsEmpty)
	}

	if e.client == nil {
		return NilClientError
	}

	getService := e.client.Get().Index(q.Index).Type(q.Index).Id(q.ID)
	if len(q.Source) > 0 {
		getService = getService.FetchSourceContext(
			eslib.NewFetchSourceContext(true).Include(q.Source...),
		)
	}
	result, err := getService.Do(ctx)

	// not found result, return here first because caller need to know if the Id is not found
	if result != nil && !result.Found {
		return errors.New(NotFound)
	}

	if err != nil {
		return err
	}

	return json.Unmarshal(*result.Source, &target)
}

// BulkInsert function
func (e *elasticV5) BulkInsert(ctx context.Context, requestParams []elastic.InsertParam) error {
	err := e.Connect()
	if err != nil {
		return err
	}
	if e.client == nil {
		return NilClientError
	}
	bulkRequest := e.client.Bulk()

	for _, requestParam := range requestParams {
		request := eslib.
			NewBulkIndexRequest().
			Index(requestParam.Index).
			Type(requestParam.Type).
			Id(requestParam.ID).
			Doc(requestParam.Doc)
		bulkRequest = bulkRequest.Add(request)
	}

	if _, err := bulkRequest.Do(ctx); err != nil {
		return err
	}
	return nil
}

// Upsert function
func (e *elasticV5) Upsert(ctx context.Context, requestParam elastic.InsertParam) (elastic.Response, error) {

	result := elastic.Response{}
	err := e.Connect()
	if err != nil {
		return result, err
	}
	if e.client == nil {
		return result, NilClientError
	}
	res, err := e.client.Index().
		Index(requestParam.Index).
		BodyJson(requestParam.Doc).
		Type(requestParam.Type).
		Id(requestParam.ID).
		Do(context.Background())

	if err != nil {
		return result, err
	}

	result.ID = res.Id
	result.Index = res.Index
	result.Result = res.Result
	result.Status = res.Status
	result.Type = res.Type
	return result, err
}

func (e *elasticV5) GetBucket(ctx context.Context, param elastic.GetBucketParam) (*eslib.AggregationBucketKeyItems, error) {
	var result *eslib.AggregationBucketKeyItems
	err := e.Connect()
	if err != nil {
		return result, err
	}
	if e.client == nil {
		return result, NilClientError
	}
	res, err := e.client.Search().Index(param.Index).Type(param.Type).Source(param.Source).Do(ctx)
	if err != nil {
		return result, err
	}

	result, found := res.Aggregations.Terms(param.Terms)
	if !found {
		return result, errors.New(NotFound)
	}

	return result, nil
}

func (e *elasticV5) Refresh(index string) {
	err := e.Connect()
	if err != nil {
		return
	}
	if e.client == nil {
		return
	}
	e.client.Refresh(index)
}

func (e *elasticV5) Delete(ctx context.Context, param elastic.DeleteParam) error {
	err := e.Connect()
	if err != nil {
		return err
	}
	if e.client == nil {
		return NilClientError
	}

	_, err = e.client.Delete().Index(param.Index).Type(param.Type).Id(param.ID).Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

// SearchAggregation : Get only aggregation results from a search query.
func (e *elasticV5) Aggregation(ctx context.Context, param elastic.SearchParam) (map[string]eslib.AggregationBucketKeyItems, error) {
	const funcName = "Aggregation"

	if param.Index == "" || param.Request == nil {
		return nil, errors.New("Invalid Search parameter")
	}
	if err := e.Connect(); err != nil {
		return nil, err
	}
	if e.client == nil {
		return nil, NilClientError
	}

	searchSvc := e.client.Search().Index(param.Index).Type(param.Index)
	searchRes, err := searchSvc.Source(param.Request).Do(ctx)
	if err != nil {
		return nil, err
	}

	// Multiple aggregations are possible.
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
		// No need to store raw data.
		agg.Aggregations = nil
		for idx := range agg.Buckets {
			bucket := agg.Buckets[idx]
			if bucket == nil {
				continue
			}
			// No need to store raw data.
			bucket.Aggregations = nil
		}
		aggregations[aggKey] = agg
	}

	return aggregations, nil
}
