package sendgrid

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/shipply-io/shipply-io-backend/util"
)

var sendgridClient *sendgrid.Client

func Init() {
	sendgridClient = sendgrid.NewSendClient(util.SendgridAPIKey)
}

func SendPasswordResetEmail(toName string, toEmail string, token string) error {

	// Create the email
	m := mail.NewV3Mail()
	m.SetFrom(mail.NewEmail("Warehance", util.SendgridFromEmail))
	m.SetTemplateID(util.SendgridResetPasswordTemplateID)

	// Add The Recipient to the personalization
	p := mail.NewPersonalization()
	p.AddTos(mail.NewEmail(toName, toEmail))

	// Add the dynamic template data to the personalization
	p.SetDynamicTemplateData("reset_password_link", util.FrontEndBaseURL+"/auth/forgot-password?token="+token)

	// Add the personalization to the email
	m.AddPersonalizations(p)

	// Send the email
	_, err := sendgridClient.Send(m)
	if err != nil {
		return err
	}

	return nil

}
