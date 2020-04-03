package gmail

import (
	"net/http"

	"google.golang.org/api/gmail/v1"
)

type gmailModule struct {
	JSONConfig string
	Token      string

	googleClient *http.Client
	service      *gmail.Service
}
