package environment

import (
	"log"
	"os"
)

const (
	ENV_OUTLOOK_email    = "API_OUTLOOK_EMAIL"
	ENV_OUTLOOK_password = "API_OUTLOOK_PASSWORD"
	ENV_OUTLOOK_addresse = "API_OUTLOOK_ADDRESSE"
)

type OutlookConfig struct {
	Email    string
	Password string
	Addresse string
}

// GetAll read all needed enviornment variables
func (o *OutlookConfig) GetAll() {
	// --------- Get Email
	if o.Email = os.Getenv(ENV_OUTLOOK_email); len(o.Email) == 0 {
		log.Fatalf("%s %s", ENV_OUTLOOK_email, environement_variable_missing)
	}

	// --------- Get Password
	if o.Password = os.Getenv(ENV_OUTLOOK_password); len(o.Password) == 0 {
		log.Fatalf("%s %s", ENV_OUTLOOK_password, environement_variable_missing)
	}

	// --------- Get Addresse
	if o.Addresse = os.Getenv(ENV_OUTLOOK_addresse); len(o.Addresse) == 0 {
		log.Fatalf("%s %s", ENV_OUTLOOK_addresse, environement_variable_missing)
	}
}
