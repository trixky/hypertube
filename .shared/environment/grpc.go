package environment

import (
	"log"
)

const (
	ENV_grpc_port         = "API_AUTH_GRPC_PORT"
	ENV_grpc_gateway_port = "API_AUTH_GRPC_GATEWAY_PORT"
)

type env_grpc struct {
	Port        int
	GatewayPort int
}

// GetAll read all needed enviornment variables
func (e *env_grpc) GetAll() {
	// --------- Get Port
	if grpc_port, err := read_port(ENV_grpc_port); err != nil {
		log.Fatal(err)
	} else {
		e.Port = grpc_port
	}

	// --------- Get GatewayPort
	if http_port, err := read_port(ENV_grpc_gateway_port); err != nil {
		log.Fatal(err)
	} else {
		e.GatewayPort = http_port
	}
}

var Grpc = env_grpc{}
