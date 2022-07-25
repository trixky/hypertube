package environment

import (
	"log"
	"os"
)

const (
	ENV_tmdb_api_key = "TMDB_API_KEY"
)

type env_tmdb struct {
	ApiKey string
}

// GetAll read all needed enviornment variables
func (e *env_tmdb) GetAll() {
	// --------- Get API key
	if e.ApiKey = os.Getenv(ENV_tmdb_api_key); len(e.ApiKey) == 0 {
		log.Fatalf("%s %s", ENV_tmdb_api_key, environement_variable_missing)
	}
}

var TMDB = env_tmdb{}
