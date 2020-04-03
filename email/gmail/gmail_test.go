package gmail_test

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	. "github.com/zainul/ark/email/gmail"
)

func TestGmailConnectError(t *testing.T) {

	// Wrong json token
	g := New("jsontoken", "gmailtoken")
	assert.EqualError(t, g.Connect(), "invalid character 'j' looking for beginning of value")

	// Wrong gmail token
	g = New(`{"installed":{"client_id":"foo","project_id":"p","auth_uri":"http://auth.com","token_uri":"http://token.com","auth_provider_x509_cert_url":"http://authcert.com","client_secret":"x","redirect_uris":["urn:ietf:wg:oauth:2.0:oob","http://localhost"]}}`, "gmailtoken")
	assert.EqualError(t, g.Connect(), "invalid character 'g' looking for beginning of value")

}

func TestGetUnreadGetListError(t *testing.T) {

	g := New(`{"installed":{"client_id":"foo","project_id":"p","auth_uri":"http://auth.com","token_uri":"http://token.com","auth_provider_x509_cert_url":"http://authcert.com","client_secret":"x","redirect_uris":["urn:ietf:wg:oauth:2.0:oob","http://localhost"]}}`, `{"access_token":"ac","token_type":"Bearer","refresh_token":"ref","expiry":"2018-11-13T18:12:24.779558896+07:00"}`)
	g.Connect()
	result, err := g.GetUnread()
	assert.Empty(t, result)
	assert.Error(t, err, "F")
}

func TestGetUnread(t *testing.T) {

	g := New(`{"installed":{"client_id":"foo","project_id":"p","auth_uri":"http://auth.com","token_uri":"http://token.com","auth_provider_x509_cert_url":"http://authcert.com","client_secret":"x","redirect_uris":["urn:ietf:wg:oauth:2.0:oob","http://localhost"]}}`, `{"access_token":"ac","token_type":"Bearer","refresh_token":"ref","expiry":"2018-11-13T18:12:24.779558896+07:00"}`)
	g.Connect()

	defer gock.Off()
	gock.New("http://token.com").Post("").Reply(200).BodyString(`{
		"access_token":"ac"
	}`)
	gock.New("https://www.googleapis.com").Get("").Reply(200).BodyString(`{
		"access_token":"ac"
	}`)

	result, err := g.GetUnread()
	assert.Empty(t, result)
	assert.Nil(t, err)
}

func TestGetUnreadWithQuery(t *testing.T) {

	g := New(`{"installed":{"client_id":"foo","project_id":"p","auth_uri":"http://auth.com","token_uri":"http://token.com","auth_provider_x509_cert_url":"http://authcert.com","client_secret":"x","redirect_uris":["urn:ietf:wg:oauth:2.0:oob","http://localhost"]}}`, `{"access_token":"ac","token_type":"Bearer","refresh_token":"ref","expiry":"2018-11-13T18:12:24.779558896+07:00"}`)
	g.Connect()

	defer gock.Off()
	gock.New("http://token.com").Post("").Reply(200).BodyString(`{
		"access_token":"ac"
	}`)
	gock.New("https://www.googleapis.com").Get("").Reply(200).BodyString(`{
		"access_token":"ac"
	}`)

	result, err := g.GetUnreadWithQuery("test", 0)
	assert.Empty(t, result)
	assert.Nil(t, err)
}
