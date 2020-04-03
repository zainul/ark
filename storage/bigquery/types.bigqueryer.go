package bigquery

import (
	bq "cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
)

type BigQueryer interface {
	Connect() error
	Query(queryParam string) Iterator
}

type bigQueryModule struct {
	JSONConfig string
	Client     *bq.Client
}

type serviceAccount struct {
	Type              string `json:"type"`
	ProjectID         string `json:"project_id"`
	PrivateKeyID      string `json:"private_key_id"`
	PrivateKey        string `json:"private_key"`
	ClientEmail       string `json:"client_email"`
	ClientID          string `json:"client_id"`
	AuthURI           string `json:"auth_uri"`
	TokenURI          string `json:"token_uri"`
	AuthProvider      string `json:"auth_provider_x509_cert_url"`
	ClientCertificate string `json:"client_x509_cert_url"`
}
type Iterator interface {
	Next(interface{}) error
}

type Testing struct {
	Date       civil.Date
	Count_User int
}
