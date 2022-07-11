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

type Api42 struct {
	RequestUrl     string
	GrantType      string
	ClientId       string
	ClientSecret   string
	RedirectionUri string
	RequestMe      string
}

func (a *Api42) GetAll() {
	// --------- get RequestUrl
	if a.RequestUrl = os.Getenv(ENV_42_redirection_request_url); len(a.RequestUrl) == 0 {
		log.Fatalf("%s environement variable missing", ENV_42_redirection_request_url)
	}

	// --------- get GrantType
	if a.GrantType = os.Getenv(ENV_42_redirection_grant_type); len(a.GrantType) == 0 {
		log.Fatalf("%s environement variable missing", ENV_42_redirection_grant_type)
	}

	// --------- get ClientId
	if a.ClientId = os.Getenv(ENV_42_redirection_client_id); len(a.ClientId) == 0 {
		log.Fatalf("%s environement variable missing", ENV_42_redirection_client_id)
	}

	// --------- get ClientSecret
	if a.ClientSecret = os.Getenv(ENV_42_redirection_client_secret); len(a.ClientSecret) == 0 {
		log.Fatalf("%s environement variable missing", ENV_42_redirection_client_secret)
	}

	// --------- get RedirectionUri
	if a.RedirectionUri = os.Getenv(ENV_42_redirection_redirection_uri); len(a.RedirectionUri) == 0 {
		log.Fatalf("%s environement variable missing", ENV_42_redirection_redirection_uri)
	}

	// --------- get RequestMe
	if a.RequestMe = os.Getenv(ENV_42_request_me); len(a.RequestMe) == 0 {
		log.Fatalf("%s environement variable missing", ENV_42_request_me)
	}
}
