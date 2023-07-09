package sendgrid

import (
	"context"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/shipply-io/shipply-io-backend/util"
)

func SendPasswordResetEmail(ctx context.Context, toName string, toEmail string, token string) error {
	c := sendgridClientFromContext(ctx)

	// Create the email
	m := mail.NewV3Mail()
	m.SetFrom(mail.NewEmail("Warehance", c.From))
	m.SetTemplateID(c.Templates["password_reset"])

	// Add The Recipient to the personalization
	p := mail.NewPersonalization()
	p.AddTos(mail.NewEmail(toName, toEmail))

	// Add the dynamic template data to the personalization
	p.SetDynamicTemplateData("reset_password_link", util.FrontendBaseURLFromContext(ctx)+"/auth/forgot-password?token="+token)

	// Add the personalization to the email
	m.AddPersonalizations(p)

	// Send the email
	_, err := c.Send(m)
	if err != nil {
		return err
	}

	return nil

}

type contextKey int

const (
	sendgridKey contextKey = iota
)

type SendgridClient struct {
	*sendgrid.Client

	From      string
	Templates map[string]string
}

func ContextWithSendgrindClient(ctx context.Context, apiKey string, from string, templates map[string]string) context.Context {
	sendgridClient := &SendgridClient{
		Client:    sendgrid.NewSendClient(apiKey),
		From:      from,
		Templates: templates,
	}
	return context.WithValue(ctx, sendgridKey, &sendgridClient)
}

func sendgridClientFromContext(ctx context.Context) *SendgridClient {
	if rv := ctx.Value(sendgridKey); rv != nil {
		return rv.(*SendgridClient)
	}
	return nil
}
