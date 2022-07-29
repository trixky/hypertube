package environment

import (
	"log"
	"os"
)

const (
	ENV_OUTLOOK_email              = "API_OUTLOOK_EMAIL"
	ENV_OUTLOOK_password           = "API_OUTLOOK_PASSWORD"
	ENV_OUTLOOK_addresse           = "API_OUTLOOK_ADDRESSE"
	ENV_EMAIL_redirection_addresse = "EMAIL_REDIRECTION_ADDRESSE"
)

type env_outlook struct {
	Email            string
	Password         string
	Addresse         string
	RedirectionEmail string
}

// GetAll read all needed enviornment variables
func (e *env_outlook) GetAll() {
	// --------- Get Email
	if e.Email = os.Getenv(ENV_OUTLOOK_email); len(e.Email) == 0 {
		log.Fatalf("%s %s", ENV_OUTLOOK_email, environement_variable_missing)
	}

	// --------- Get Password
	if e.Password = os.Getenv(ENV_OUTLOOK_password); len(e.Password) == 0 {
		log.Fatalf("%s %s", ENV_OUTLOOK_password, environement_variable_missing)
	}

	// --------- Get Addresse
	if e.Addresse = os.Getenv(ENV_OUTLOOK_addresse); len(e.Addresse) == 0 {
		log.Fatalf("%s %s", ENV_OUTLOOK_addresse, environement_variable_missing)
	}

	// --------- Get RedirectionEmail
	e.RedirectionEmail = os.Getenv(ENV_EMAIL_redirection_addresse)
}

var Outlook = env_outlook{}
