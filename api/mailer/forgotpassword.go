package mailer

import (
	"github.com/matcornic/hermes/v2"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/victorsteven/fullstack/api/models"
	"os"
)

func SendResetPassword(reset_password *models.ResetPassword) error {
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

	var forgotUrl string
	if os.Getenv("APP_ENV") == "production" {
		forgotUrl = "https://seamflow.com/resetpassword/" + reset_password.Token
	} else {
		forgotUrl = "http://127.0.0.1:3000/resetpassword/" + reset_password.Token
	}

	email := hermes.Email{
		Body: hermes.Body{
			Name: reset_password.Email,
			Intros: []string{
				"Welcome to SeamFlow! Good to have you here.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click this link to reset your password",
					Button: hermes.Button{
						Color: "#FFFFFF", // Optional action button color
						Text:  "Reset Password",
						Link: forgotUrl,
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
	subject := "Reset Password"
	to := mail.NewEmail("Reset Password", reset_password.Email)
	message := mail.NewSingleEmail(from, subject, to, emailBody, emailBody)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err = client.Send(message)
	if err != nil {
		return err
	}
	return nil
}
