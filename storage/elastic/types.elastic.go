package elastic

import (
	"context"
	"encoding/json"

	eslib "gopkg.in/olivere/elastic.v5"
)

// An Elastic offers a standard interface for connect to elasticsearch service
type Elastic interface {
	Connect() error
	GetClient() interface{}
	Search(context.Context, SearchParam) ([]string, []json.RawMessage, error)
	Suggest(context.Context, SuggestParam) ([]string, []json.RawMessage, error)
	Get(ctx context.Context, q GetParam, target interface{}) error
	BulkInsert(ctx context.Context, requestParams []InsertParam) error
	Upsert(ctx context.Context, requestParam InsertParam) (Response, error)
	GetBucket(ctx context.Context, param GetBucketParam) (*eslib.AggregationBucketKeyItems, error)
	Delete(ctx context.Context, param DeleteParam) error
	Refresh(index string)
	Aggregation(ctx context.Context, param SearchParam) (map[string]eslib.AggregationBucketKeyItems, error)
}

type DeleteParam struct {
	Index string
	Type  string
	ID    string
}

// SearchParam for elastic search query
type SearchParam struct {
	Index   string
	Request interface{}
}

// SuggestParam for elastic suggest query
type SuggestParam struct {
	SuggestName string
	Index       string
	Request     interface{}
}

// GetParam for elastic get query
type GetParam struct {
	Index  string
	ID     string
	Source []string
}

// InsertParam for elastic push request
type InsertParam struct {
	Index string
	Type  string
	ID    string
	Doc   interface{}
}

// Response of elastic
type Response struct {
	Index  string `json:"_index,omitempty"`
	Type   string `json:"_type,omitempty"`
	ID     string `json:"_id,omitempty"`
	Result string `json:"result,omitempty"`
	Status int    `json:"status,omitempty"`
}

type GetBucketParam struct {
	Index  string
	Type   string
	Source interface{}
	Terms  string
}
