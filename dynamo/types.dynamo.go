package dynamo

import "github.com/guregu/dynamo"

type DynamoCfg struct {
	TableName  string
	Region     string
	Endpoint   string
	DisableSSL bool
}

/*DynStore : model for dynamo datastore*/
type DynStore struct {
	DB    *dynamo.DB
	Table string
}
