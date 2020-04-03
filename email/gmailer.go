package email

// Gmailer interface function
type Gmailer interface {
	Connect() error
	GetUnread() ([]*GmailContent, error)
	GetUnreadWithQuery(query string, maxResults int64) ([]*GmailContent, error)
	UpdateLabels(string, []string, []string) error
}

// GmailContent struct
type GmailContent struct {
	ID        string
	From      string
	Subject   string
	Timestamp int64
	Body      string
}
