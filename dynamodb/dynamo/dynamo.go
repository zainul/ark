package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/guregu/dynamo"
	"github.com/zainul/ark/dynamodb"
)

type (
	// Config : DynamoDB config parameters.
	Config struct {
		DefaultTable string
		Region       string
		Endpoint     string
		SecretKey    string
		AccessKey    string
		Token        string
		DisableSSL   bool
	}

	// Dynamo : DynamoDB struct implementation.
	Dynamo struct {
		db           *wrappedDB
		defaultTable string
	}

	// wrappedDB : Wrap the DynamoDB to make it possible to use inteface and mockable.
	wrappedDB struct {
		db *dynamo.DB
	}
)

// NewDynamo : Create new DynamoDB connection.
func NewDynamo(config Config) *Dynamo {
	return &Dynamo{
		defaultTable: config.DefaultTable,
		db: &wrappedDB{
			db: dynamo.New(
				session.New(),
				aws.NewConfig().
					WithRegion(config.Region).
					WithEndpoint(config.Endpoint).
					WithCredentials(credentials.NewStaticCredentials(
						config.AccessKey,
						config.SecretKey,
						config.Token,
					)),
			),
		},
	}
}

// GetDB : Get DynamoDB client.
func (dyn *Dynamo) GetDB() dynamodb.DynamoClient {
	return dyn.db
}

// GetDefaultTable : Get default table name used by the DynamoDB client.
func (dyn *Dynamo) GetDefaultTable() string {
	return dyn.defaultTable
}

func (wdb *wrappedDB) Table(name string) dynamodb.DynamoTable {
	return wdb.db.Table(name)
}

func (wdb *wrappedDB) Client() dynamodbiface.DynamoDBAPI {
	return wdb.db.Client()
}

func (wdb *wrappedDB) ListTables() *dynamo.ListTables {
	return wdb.db.ListTables()
}

func (wdb *wrappedDB) CreateTable(name string, from interface{}) *dynamo.CreateTable {
	return wdb.CreateTable(name, from)
}

func (wdb *wrappedDB) GetTx() *dynamo.GetTx {
	return wdb.db.GetTx()
}

func (wdb *wrappedDB) WriteTx() *dynamo.WriteTx {
	return wdb.db.WriteTx()
}
