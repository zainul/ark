package dynamo_test

import (
	"testing"

	"github.com/zainul/ark/dynamo"
)

func TestSetupDynamo(t *testing.T) {
	dynamo.SetupDynamo(dynamo.DynamoCfg{})
}

func TestNewDynamoStorage(t *testing.T) {
	dynamo.NewDynamoStorage(nil, "")
}
