package environment

import (
	"log"
	"os"
)

const (
	ENV_demo_mode = "DEMO_MODE"
)

type env_demo struct {
	Demo bool
}

// GetAll read all needed enviornment variables
func (e *env_demo) GetAll() {
	// --------- Get Domain
	if demo := os.Getenv(ENV_demo_mode); len(demo) == 0 {
		log.Fatalf("%s %s", ENV_demo_mode, environement_variable_missing)
	} else if demo == "true" {
		e.Demo = true
	} else if demo == "false" {
		e.Demo = false
	} else {
		log.Fatalf("%s %s", ENV_demo_mode, environement_variable_corrupted_boolean)
	}
}

var Demo = env_demo{}
