package messaging

import (
	"canvas/model"
	"context"
	"embed"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	marketingMessageStream   = "broadcast"
	transactionMessageStream = "outbound"
)

// nameAndEmail combo, of the form "Name <email@example.com>"
type nameAndEmail = string

//go:embed emails
var emails embed.FS

type Emailer struct {
	baseURL           string
	client            *http.Client
	log               *zap.Logger
	marketingFrom     nameAndEmail
	token             string
	transactionalFrom nameAndEmail
}

type NewEmailerOptions struct {
	BaseURL                   string
	Log                       *zap.Logger
	MarketingEmailAddress     string
	MarketingEmailName        string
	TransactionalEmailAddress string
	TransactionalEmailName    string
	Token                     string
}

func NewEmailer(opts NewEmailerOptions) *Emailer {
	return &Emailer{
		baseURL:           opts.BaseURL,
		client:            &http.Client{Timeout: 3 * time.Second},
		log:               opts.Log,
		marketingFrom:     createNameAndEmail(opts.MarketingEmailName, opts.MarketingEmailAddress),
		token:             opts.Token,
		transactionalFrom: createNameAndEmail(opts.TransactionalEmailName, opts.TransactionalEmailAddress),
	}
}

// createNameAndEmail returns a name and email string ready for inserting into From and To fields.
func createNameAndEmail(name, email string) nameAndEmail {
	return fmt.Sprintf("%v <%v>", name, email)
}

func (e *Emailer) SendNewsletterConfirmationEmail(ctx context.Context, to model.Email, token string) {
	keywords := map[string]string{
		"base_url": e.baseURL,
		"action_url": e.baseURL + "/newsletter/confirm?token=" + token
	}

	return e.sed
}
