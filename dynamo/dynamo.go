package dynamo

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

var UserLogin *DynStore

func SetupDynamo(cfg DynamoCfg) {
	if UserLogin != nil {
		return
	}

	UserLogin = NewDynamoStorage(
		aws.NewConfig().WithRegion(cfg.Region).WithEndpoint(cfg.Endpoint),
		cfg.TableName,
	)

}

/*NewDynamoStorage : this is an implementation of Storage that uses dynamodb*/
func NewDynamoStorage(cfg *aws.Config, tablename string) *DynStore {
	if cfg == nil {
		log.Println("Error EMPTY CONFIG")
		cfg = &aws.Config{Region: aws.String("ap-southeast-1"), Endpoint: aws.String("http://127.0.0.1:4567")}
	}

	db := dynamo.New(session.New(), cfg)

	return &DynStore{
		DB:    db,
		Table: tablename,
	}
}
