package mailer

import (
	"github.com/matcornic/hermes/v2"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/victorsteven/forum/api/models"
	"os"
)

func SendMail(user *models.User) error {
	// Configure hermes by setting a theme and your product info
	h := hermes.Hermes{
		// Optional Theme
		// Theme: new(Default)
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: "SeamFlow",
			Link: "https://seamflow.com",
			// Optional product logo
			//Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}
	email := hermes.Email{
		Body: hermes.Body{
			Name: user.Username,
			Intros: []string{
				"Welcome to SeamFlow! Good to have you here.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with SeamFlow, please click here:",
					Button: hermes.Button{
						Color: "#FFFFFF", // Optional action button color
						Text:  "Confirm your account",
						Link:  "https://hermes-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}
	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		return err
	}
	from := mail.NewEmail("SeamFlow", os.Getenv("SENDGRID_FROM"))
	subject := "Registration Successful"
	to := mail.NewEmail("Welcome", user.Email)
	message := mail.NewSingleEmail(from, subject, to, emailBody, emailBody)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err = client.Send(message)
	if err != nil {
		return err
	}
	return nil
}
