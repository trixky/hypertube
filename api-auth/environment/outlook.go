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

func (o *OutlookConfig) GetAll() {
	// --------- get Email
	if o.Email = os.Getenv(ENV_OUTLOOK_email); len(o.Email) == 0 {
		log.Fatalf("%s environement variable missing", ENV_OUTLOOK_email)
	}

	// --------- get Password
	if o.Password = os.Getenv(ENV_OUTLOOK_password); len(o.Password) == 0 {
		log.Fatalf("%s environement variable missing", ENV_OUTLOOK_password)
	}

	// --------- get Addresse
	if o.Addresse = os.Getenv(ENV_OUTLOOK_addresse); len(o.Addresse) == 0 {
		log.Fatalf("%s environement variable missing", ENV_OUTLOOK_addresse)
	}
}
