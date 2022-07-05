package environment

import (
	"log"
	"os"
)

const (
	ENV_42_redirection_request_url     = "42_REDIRECTION_REQUEST_URL"
	ENV_42_redirection_grant_type      = "42_REDIRECTION_GRANT_TYPE"
	ENV_42_redirection_client_id       = "42_REDIRECTION_CLIENT_ID"
	ENV_42_redirection_client_secret   = "42_REDIRECTION_CLIENT_SECRET"
	ENV_42_redirection_redirection_uri = "42_REDIRECTION_REDIRECT_URI_ID"
	ENV_42_request_me                  = "42_REQUEST_ME"
)

type Api42 struct {
	RequestUrl     string
	GrantType      string
	ClientId       string
	ClientSecret   string
	RedirectionUri string
	RequestMe      string
}

func (r *Api42) GetAll() {
	// --------- get RequestUrl
	if r.RequestUrl = os.Getenv(ENV_42_redirection_request_url); len(r.RequestUrl) == 0 {
		log.Fatalf("%s environement variable missing", ENV_42_redirection_request_url)
	}

	// --------- get GrantType
	if r.GrantType = os.Getenv(ENV_42_redirection_grant_type); len(r.GrantType) == 0 {
		log.Fatalf("%s environement variable missing", ENV_42_redirection_grant_type)
	}

	// --------- get ClientId
	if r.ClientId = os.Getenv(ENV_42_redirection_client_id); len(r.ClientId) == 0 {
		log.Fatalf("%s environement variable missing", ENV_42_redirection_client_id)
	}

	// --------- get ClientSecret
	if r.ClientSecret = os.Getenv(ENV_42_redirection_client_secret); len(r.ClientSecret) == 0 {
		log.Fatalf("%s environement variable missing", ENV_42_redirection_client_secret)
	}

	// --------- get RedirectionUri
	if r.RedirectionUri = os.Getenv(ENV_42_redirection_redirection_uri); len(r.RedirectionUri) == 0 {
		log.Fatalf("%s environement variable missing", ENV_42_redirection_redirection_uri)
	}

	// --------- get RequestMe
	if r.RequestMe = os.Getenv(ENV_42_request_me); len(r.RequestMe) == 0 {
		log.Fatalf("%s environement variable missing", ENV_42_request_me)
	}
}
