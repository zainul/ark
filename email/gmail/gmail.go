package gmail

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"html"
	"log"
	"strings"

	"github.com/grokify/html-strip-tags-go"
	"github.com/zainul/ark/email"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmailAPI "google.golang.org/api/gmail/v1"
)

var replacer = strings.NewReplacer("-", "+", "_", "/")

// New gmail module
func New(jsonConfig, gmailToken string) email.Gmailer {

	return &gmailModule{
		JSONConfig: jsonConfig,
		Token:      gmailToken,
	}
}

// Connect to oauth
func (g *gmailModule) Connect() error {

	// Initialize Google Config
	googleConfig, err := google.ConfigFromJSON([]byte(g.JSONConfig), gmailAPI.GmailReadonlyScope)
	if err != nil {
		return err
	}

	// Initialize Google Client
	token := &oauth2.Token{}
	if err = json.Unmarshal([]byte(g.Token), token); err != nil {
		return err
	}

	g.googleClient = googleConfig.Client(context.Background(), token)

	// Create gmail service
	g.service, _ = gmailAPI.New(g.googleClient)

	return nil
}

// GetUnread -- latest 50 emails
func (g *gmailModule) GetUnread() ([]*email.GmailContent, error) {

	// Get unread from inbox -- last 50 emails
	response, err := g.service.Users.Messages.List("me").LabelIds("INBOX", "UNREAD").MaxResults(50).Do()
	if err != nil {
		return nil, err
	}

	return g.processMessage(response.Messages)

}

// GetUnreadWithQuery func
func (g *gmailModule) GetUnreadWithQuery(query string, maxResults int64) ([]*email.GmailContent, error) {

	if maxResults == 0 {
		maxResults = 50
	}

	// Get unread from inbox
	response, err := g.service.Users.Messages.List("me").Q(query).MaxResults(maxResults).Do()
	if err != nil {
		return nil, err
	}

	return g.processMessage(response.Messages)

}

func (g *gmailModule) processMessage(messages []*gmailAPI.Message) ([]*email.GmailContent, error) {

	var contents = []*email.GmailContent{}
	// Process messages
	totalMsg := len(messages)
	for m := totalMsg - 1; m > 0; m-- {

		if message, err := g.service.Users.Messages.Get("me", messages[m].Id).Do(); err != nil {

			// Fail to get email details
			continue

		} else if message.Payload != nil {

			content := &email.GmailContent{
				Timestamp: message.InternalDate / 1000,
				ID:        message.Id,
			}

			// Get From and Subject
			for h := range message.Payload.Headers {
				if message.Payload.Headers[h].Name == "From" {
					content.From = message.Payload.Headers[h].Value
				} else if message.Payload.Headers[h].Name == "Subject" {
					content.Subject = message.Payload.Headers[h].Value
				}
			}

			var decoded string

			// Message has no part -- get content from body
			if message.Payload.Body != nil && len(message.Payload.Parts) <= 1 {
				decoded += getBody(message.Payload.Body.Data)
			}

			// Multipart messages
			for p := range message.Payload.Parts {

				// Parts has multi inner-parts
				if message.Payload.Parts[p].Parts != nil {
					for prs := range message.Payload.Parts[p].Parts {
						decoded += getBody(message.Payload.Parts[p].Parts[prs].Body.Data)
					}

				} else {
					decoded += getBody(message.Payload.Parts[p].Body.Data)

				}
			}

			content.Body = decoded
			contents = append(contents, content)

		}

	}

	return contents, nil
}

// UpdateLabels by message id
func (g *gmailModule) UpdateLabels(id string, addLabels, removeLabels []string) error {

	modifyParam := &gmailAPI.ModifyMessageRequest{
		AddLabelIds:    addLabels,
		RemoveLabelIds: removeLabels,
	}

	// Hit Gmail API
	_, err := g.service.Users.Messages.Modify("me", id, modifyParam).Do()
	return err
}

func getBody(data string) string {

	// Decode base64 to string
	encodedData := replacer.Replace(data)
	dc, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		log.Println("func getBody", err)
		return ""
	}

	// Unescape HTML string and remove all tags
	body := html.UnescapeString(strip.StripTags(string(dc)))

	// Remove redundant spaces
	body = strings.Join(strings.Fields(body), " ")

	// Replace HTML Header
	body = strings.Replace(body, `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">`, "", -1)

	// Replace " with \" to prevent broken JSON
	body = strings.Replace(body, `"`, `\"`, -1)

	return body

}
