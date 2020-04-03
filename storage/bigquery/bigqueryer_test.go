package bigquery_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/zainul/ark/storage/bigquery"
)

// Skip this unit test because breaking dependency, eventually this packages will be removed due to unused
func TestNew(t *testing.T) {

	jsonConfig := "xxx"
	bigQuery := New(jsonConfig)
	assert.NotNil(t, bigQuery)

	connection := bigQuery.Connect()
	assert.Error(t, connection)

	bigQuery = New(`{"type":"service_account","project_id":"xxx-970","private_key_id":"xxx","private_key":"xxxx","client_email":"xxxx","client_id":"xxx","auth_uri":"x","token_uri":"xxxx","auth_provider_x509_cert_url":"xxxx","client_x509_cert_url":"xxxx"}`)
	connection = bigQuery.Connect()

	assert.Nil(t, connection)
}
