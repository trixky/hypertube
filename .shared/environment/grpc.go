package environment

import (
	"log"
)

type env_grpc struct {
	Port        int
	GatewayPort int
}

// GetAll read all needed enviornment variables
func (e *env_grpc) GetAll(config *Config) {
	// --------- Get Port
	if grpc_port, err := ReadPort(config.ENV_grpc_port); err != nil {
		log.Fatal(err)
	} else {
		e.Port = grpc_port
	}

	// --------- Get GatewayPort
	if grpc_gateway_port, err := ReadPort(config.ENV_grpc_gateway_port); err != nil {
		log.Fatal(err)
	} else {
		e.GatewayPort = grpc_gateway_port
	}
}

var Grpc = env_grpc{}
