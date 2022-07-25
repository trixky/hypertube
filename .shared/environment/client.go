package environment

import (
	"log"
	"os"
)

const (
	ENV_domain = "DOMAIN"
	ENV_port   = "CLIENT_PORT"
)

type env_client struct {
	Domain string
	Port   int
}

// GetAll read all needed enviornment variables
func (e *env_client) GetAll() {
	// --------- Get Domain
	if e.Domain = os.Getenv(ENV_domain); len(e.Domain) == 0 {
		log.Fatalf("%s %s", ENV_domain, environement_variable_missing)
	}

	// --------- Get Client port
	if port, err := ReadPort(ENV_port); err != nil {
		log.Fatal(err)
	} else {
		e.Port = port
	}
}

var Client = env_client{}
