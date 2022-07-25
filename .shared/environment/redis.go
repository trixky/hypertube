package environment

import (
	"log"
	"os"
)

const (
	ENV_redis_host = "REDIS_HOST"
	ENV_redis_port = "REDIS_PORT"
)

type env_redis struct {
	Host string
	Port string
}

// GetAll read all needed enviornment variables
func (e *env_redis) GetAll() {
	// --------- Get Host
	if e.Host = os.Getenv(ENV_redis_host); len(e.Host) == 0 {
		log.Fatalf("%s %s", ENV_redis_host, environement_variable_missing)
	}

	// --------- Get Port
	if e.Port = os.Getenv(ENV_redis_port); len(e.Port) == 0 {
		log.Fatalf("%s %s", ENV_redis_port, environement_variable_missing)
	}
}

var Redis = env_redis{}
