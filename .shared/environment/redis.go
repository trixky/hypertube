package environment

import (
	"log"
	"os"
)

const (
	ENV_redis_host = "REDIS_HOST"
)

type env_redis struct {
	RedisHost string
}

// GetAll read all needed enviornment variables
func (e *env_redis) GetAll() {
	// --------- Get RedisHost
	if e.RedisHost = os.Getenv(ENV_redis_host); len(e.RedisHost) == 0 {
		log.Fatalf("%s %s", ENV_redis_host, environement_variable_missing)
	}
}

var Redis = env_redis{}
