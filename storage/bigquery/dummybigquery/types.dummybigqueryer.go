package dummybigquery

import "github.com/zainul/ark/storage/bigquery"

type DummyBigQueryer interface {
	Connect() error
	Query(queryParam string) bigquery.Iterator
}

type IteratorDummy struct {
	ExpectedResult []string
	CurrentIndex   int
}

type dummybigQueryModule struct {
	Parameters DummyArgs
}

type DummyArgs struct {
	ExpectedResult    []string
	IsConnectionError bool
	CurrentIndex      int
}

type Target struct {
	Count int `json:"count"`
	User  int `json:"user"`
}
