// https://www.courier.com/guides/golang-send-email/

package email

import (
	"net/smtp"

	"github.com/trixky/hypertube/api-auth/environment"
)

func SendRegistrationConfirmation(user_email string) error {
	subject := "Registrion confirmation to Hypertube"
	body := "We confirm your registration on our platform! Welcome and good viewing"
	msg := []byte("From: " + environment.E.OUTLOOKConfig.Email + "\r\n" +
		"To: " + user_email + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n")

	auth := loginAuthOutlook(environment.E.OUTLOOKConfig.Email, environment.E.OUTLOOKConfig.Password)
	return smtp.SendMail(environment.E.OUTLOOKConfig.Addresse, auth, environment.E.OUTLOOKConfig.Email, []string{user_email}, msg)
}

func SendTokenPassword(user_email string, token_password string) error {
	subject := "Recover your Hypertube password"
	redirect_url := "http://localhost:4040/recover/apply" + "?password_token=" + token_password
	body := "You asked to recover your password.\r\nDon't wait, the link is temporary!\r\n\n" + redirect_url + "\r\n"

	msg := []byte("From: " + environment.E.OUTLOOKConfig.Email + "\r\n" +
		"To: " + user_email + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n")

	auth := loginAuthOutlook(environment.E.OUTLOOKConfig.Email, environment.E.OUTLOOKConfig.Password)
	return smtp.SendMail(environment.E.OUTLOOKConfig.Addresse, auth, environment.E.OUTLOOKConfig.Email, []string{user_email}, msg)
}
