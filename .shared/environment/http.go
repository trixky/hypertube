package environment

import (
	"log"
)

const (
	ENV_http_port = "API_AUTH_HTTP_PORT"
)

type env_http struct {
	HttpPort int
}

// GetAll read all needed enviornment variables
func (e *env_http) GetAll() {
	// --------- Get HttpPort
	if redirect_port, err := read_port(ENV_http_port); err != nil {
		log.Fatal(err)
	} else {
		e.HttpPort = redirect_port
	}
}

var Http = env_http{}
