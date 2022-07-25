package environment

import (
	"log"
)

const ()

type env_http struct {
	Port int
}

// GetAll read all needed enviornment variables
func (e *env_http) GetAll(config *Config) {
	// --------- Get Port
	if redirect_port, err := ReadPort(config.ENV_http_port); err != nil {
		log.Fatal(err)
	} else {
		e.Port = redirect_port
	}
}

var Http = env_http{}
