package messaging

import (
	"bytes"
	"canvas/model"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func (e *Emailer) SendNewsletterConfirmationEmail(ctx context.Context, to model.Email, token string) error {
	keywords := map[string]string{
		"base_url":   e.baseURL,
		"action_url": e.baseURL + "/newsletter/confirm?token=" + token,
	}

	// this will return error/nil if anything goes wrong with send email
	return e.send(
		ctx,
		requestBody{
			MessageStream: transactionMessageStream,
			From:          e.transactionalFrom,
			To:            to.String(),
			Subject:       "AYO confirm your subscription to canvas newsletter",
			HtmlBody:      getEmail("confirmation_email.html", keywords),
			TextBody:      getEmail("confirmation_email.txt", keywords),
		},
	)
}

// requestBody used in Emailer.send
// See https://postmarkapp.com/developer/user-guide/send-email-with-api

type requestBody struct {
	MessageStream string
	From          nameAndEmail
	To            nameAndEmail
	Subject       string
	HtmlBody      string
	TextBody      string
}

// send using postmark API
func (e *Emailer) send(ctx context.Context, body requestBody) error {
	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshaling request body to json: %w", err)
	}

	// create NEW request
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.postmarkapp.com/email",
		bytes.NewReader(bodyAsBytes))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	// set headers
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Postmark-Server-Token", e.token)

	// make POST call with data
	response, err := e.client.Do(request)
	if err != nil {
		return fmt.Errorf("error marking request: %w", err)
	}
	defer response.Body.Close()

	// read response
	bodyAsBytes, err = io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	// check statusCode
	if response.StatusCode > 299 {
		e.log.Info("Error sending email",
			zap.Int("status", response.StatusCode), zap.String("response", string(bodyAsBytes)))
		return fmt.Errorf("error sending email, got status: %d", response.StatusCode)
	}

	// everything good, return no error
	return nil
}

// getEmail from given path, panic on errors (because this is important to bundle with codebase)
func getEmail(path string, keywords map[string]string) string {
	email, err := emails.ReadFile("emails/" + path)
	if err != nil {
		panic(err.Error())
	}

	emailString := string(email)
	for keyword, replacement := range keywords {
		emailString = strings.ReplaceAll(emailString, "{{"+keyword+"}}", replacement)
	}

	return emailString
}
