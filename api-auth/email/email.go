package email

// https://www.courier.com/guides/golang-send-email/

import (
	"net/smtp"

	"github.com/trixky/hypertube/.shared/environment"
)

// SendRegistrationConfirmation sends a registration confirmation email after an account creation
func SendRegistrationConfirmation(user_email string) error {
	// Generates the subject of the email
	subject := "Registrion confirmation to Hypertube"
	// Generates the body of the email
	body := "We confirm your registration on our platform! Welcome and good viewing"
	// Generates the msg of the email from all these parts
	msg := []byte("From: " + environment.Outlook.Email + "\r\n" +
		"To: " + user_email + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n")

	// Generates the authentification part of the request
	auth := loginAuthOutlook(environment.Outlook.Email, environment.Outlook.Password)

	// Send the email
	err := smtp.SendMail(environment.Outlook.Addresse, auth, environment.Outlook.Email, []string{user_email}, msg)

	return err
}

// SendPasswordToken sends a password token email
func SendPasswordToken(user_email string, password_token string) error {
	// Generates the subject of the email
	subject := "Recover your Hypertube password"
	// Generates the redirect_url of the email
	redirect_url := "http://localhost:4040/recover/apply" + "?password_token=" + password_token
	// Generates the body of the email
	body := "You asked to recover your password.\r\nDon't wait, the link is temporary!\r\n\n" + redirect_url + "\r\n"
	// Generates the msg of the email from all these parts
	msg := []byte("From: " + environment.Outlook.Email + "\r\n" +
		"To: " + user_email + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n")

	// Generates the authentification part of the request
	auth := loginAuthOutlook(environment.Outlook.Email, environment.Outlook.Password)

	// Send the email
	err := smtp.SendMail(environment.Outlook.Addresse, auth, environment.Outlook.Email, []string{user_email}, msg)

	return err
}
