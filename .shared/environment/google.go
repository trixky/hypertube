package environment

import (
	"log"
	"os"
)

const (
	ENV_GOOGLE_client_id      = "API_GOOGLE_client_id"
	ENV_GOOGLE_client_sercret = "API_GOOGLE_client_secret"
	ENV_GOOGLE_redirect_url   = "API_GOOGLE_redirect_url"
	ENV_GOOGLE_scope_email    = "API_GOOGLE_scope_email"
	ENV_GOOGLE_scope_userinfo = "API_GOOGLE_scope_userinfo"
	ENV_GOOGLE_userinfo_url   = "API_GOOGLE_userinfo_url"
)

type env_api_google struct {
	ClientId       string
	ClientSecret   string
	RedirectURL    string
	ScopeEmail     string
	ScopesUserinfo string
	UserInfoURL    string
}

// GetAll read all needed enviornment variables
func (e *env_api_google) GetAll() {
	// --------- Get ClientId
	if e.ClientId = os.Getenv(ENV_GOOGLE_client_id); len(e.ClientId) == 0 {
		log.Fatalf("%s %s", ENV_GOOGLE_client_id, environement_variable_missing)
	}

	// --------- Get ClientSecret
	if e.ClientSecret = os.Getenv(ENV_GOOGLE_client_sercret); len(e.ClientSecret) == 0 {
		log.Fatalf("%s %s", ENV_GOOGLE_client_sercret, environement_variable_missing)
	}

	// --------- Get RedirectURL
	if e.RedirectURL = os.Getenv(ENV_GOOGLE_redirect_url); len(e.RedirectURL) == 0 {
		log.Fatalf("%s %s", ENV_GOOGLE_redirect_url, environement_variable_missing)
	}

	// --------- Get ScopeEmail
	if e.ScopeEmail = os.Getenv(ENV_GOOGLE_scope_email); len(e.ScopeEmail) == 0 {
		log.Fatalf("%s %s", ENV_GOOGLE_scope_email, environement_variable_missing)
	}

	// --------- Get ScopesUserinfo
	if e.ScopesUserinfo = os.Getenv(ENV_GOOGLE_scope_userinfo); len(e.ScopesUserinfo) == 0 {
		log.Fatalf("%s %s", ENV_GOOGLE_scope_userinfo, environement_variable_missing)
	}

	// --------- Get UserInfoURL
	if e.UserInfoURL = os.Getenv(ENV_GOOGLE_userinfo_url); len(e.UserInfoURL) == 0 {
		log.Fatalf("%s %s", ENV_GOOGLE_userinfo_url, environement_variable_missing)
	}
}

var ApiGoogle = env_api_google{}
