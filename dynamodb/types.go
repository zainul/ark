package dynamodb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/guregu/dynamo"
)

type (
	// Dynamo : DynamoDB container interface.
	Dynamo interface {
		GetDB() DynamoClient
		GetDefaultTable() string
	}

	// DynamoClient : DynamoDB client interface.
	DynamoClient interface {
		Client() dynamodbiface.DynamoDBAPI
		ListTables() *dynamo.ListTables
		Table(name string) DynamoTable
		CreateTable(name string, from interface{}) *dynamo.CreateTable
		GetTx() *dynamo.GetTx
		WriteTx() *dynamo.WriteTx
	}

	// DynamoTable : DynamoDB table interface.
	DynamoTable interface {
		Name() string
		Batch(hashAndRangeKeyName ...string) dynamo.Batch
		Check(hashKey string, value interface{}) *dynamo.ConditionCheck
		Delete(name string, value interface{}) *dynamo.Delete
		Describe() *dynamo.DescribeTable
		Put(item interface{}) *dynamo.Put
		Get(name string, value interface{}) *dynamo.Query
		Scan() *dynamo.Scan
		DeleteTable() *dynamo.DeleteTable
		UpdateTTL(attribute string, enabled bool) *dynamo.UpdateTTL
		DescribeTTL() *dynamo.DescribeTTL
		Update(hashKey string, value interface{}) *dynamo.Update
		UpdateTable() *dynamo.UpdateTable
	}
)
