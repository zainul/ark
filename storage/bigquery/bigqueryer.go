package bigquery

import (
	"context"
	"encoding/json"
	"log"

	bq "cloud.google.com/go/bigquery"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// New bigquery module
func New(jsonConfig string) BigQueryer {

	return &bigQueryModule{
		JSONConfig: jsonConfig,
	}
}

// Connect to oauth
func (b *bigQueryModule) Connect() error {

	creds, err := google.CredentialsFromJSON(context.Background(), []byte(b.JSONConfig), bq.Scope)
	if err != nil {
		// TODO: handle error.
		log.Println(err)
		return err
	}

	confServiceAccount := serviceAccount{}
	err = json.Unmarshal([]byte(b.JSONConfig), &confServiceAccount)

	if err != nil {
		// handle wrong json parameter
		log.Println(err)
		return err
	}

	client, err := bq.NewClient(context.Background(), confServiceAccount.ProjectID, option.WithCredentials(creds))
	if err != nil {
		// TODO: handle error.
		log.Println(err)
		return err
	}
	// Use the client.
	b.Client = client

	return nil
}

func (b *bigQueryModule) Query(queryParam string) Iterator {

	query := b.Client.Query(queryParam)
	query.UseStandardSQL = true

	result, err := query.Read(context.Background())

	if err != nil {
		log.Println(err)
		return nil
	}
	return result
}
