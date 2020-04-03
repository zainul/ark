package v5

import "errors"

// exported const
const (
	NotFound     = "Result not found"
	IndexIsEmpty = "Index is empty"
	IDIsEmpty    = "ID is empty"
)

var (
	NilClientError = errors.New("client is nil")
)
