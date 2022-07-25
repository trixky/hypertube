package environment

import (
	"log"
	"os"
)

const (
	ENV_42_redirection_request_url     = "API_42_REDIRECTION_REQUEST_URL"
	ENV_42_redirection_grant_type      = "API_42_REDIRECTION_GRANT_TYPE"
	ENV_42_redirection_client_id       = "API_42_REDIRECTION_CLIENT_ID"
	ENV_42_redirection_client_secret   = "API_42_REDIRECTION_CLIENT_SECRET"
	ENV_42_redirection_redirection_uri = "API_42_REDIRECTION_REDIRECT_URI_ID"
	ENV_42_request_me                  = "API_42_REQUEST_ME"
)

type env_api_42 struct {
	RequestUrl     string
	GrantType      string
	ClientId       string
	ClientSecret   string
	RedirectionUri string
	RequestMe      string
}

// GetAll read all needed enviornment variables
func (e *env_api_42) GetAll() {
	// --------- Get RequestUrl
	if e.RequestUrl = os.Getenv(ENV_42_redirection_request_url); len(e.RequestUrl) == 0 {
		log.Fatalf("%s %s", ENV_42_redirection_request_url, environement_variable_missing)
	}

	// --------- Get GrantType
	if e.GrantType = os.Getenv(ENV_42_redirection_grant_type); len(e.GrantType) == 0 {
		log.Fatalf("%s %s", ENV_42_redirection_grant_type, environement_variable_missing)
	}

	// --------- Get ClientId
	if e.ClientId = os.Getenv(ENV_42_redirection_client_id); len(e.ClientId) == 0 {
		log.Fatalf("%s %s", ENV_42_redirection_client_id, environement_variable_missing)
	}

	// --------- Get ClientSecret
	if e.ClientSecret = os.Getenv(ENV_42_redirection_client_secret); len(e.ClientSecret) == 0 {
		log.Fatalf("%s %s", ENV_42_redirection_client_secret, environement_variable_missing)
	}

	// --------- Get RedirectionUri
	if e.RedirectionUri = os.Getenv(ENV_42_redirection_redirection_uri); len(e.RedirectionUri) == 0 {
		log.Fatalf("%s %s", ENV_42_redirection_redirection_uri, environement_variable_missing)
	}

	// --------- Get RequestMe
	if e.RequestMe = os.Getenv(ENV_42_request_me); len(e.RequestMe) == 0 {
		log.Fatalf("%s %s", ENV_42_request_me, environement_variable_missing)
	}
}

var Api42 = env_api_42{}
