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

type ApiGoogle struct {
	ClientId       string
	ClientSecret   string
	RedirectURL    string
	ScopeEmail     string
	ScopesUserinfo string
	UserInfoURL    string
}

// GetAll read all needed enviornment variables
func (a *ApiGoogle) GetAll() {
	// --------- Get ClientId
	if a.ClientId = os.Getenv(ENV_GOOGLE_client_id); len(a.ClientId) == 0 {
		log.Fatalf("%s %s", ENV_GOOGLE_client_id, environement_variable_missing)
	}

	// --------- Get ClientSecret
	if a.ClientSecret = os.Getenv(ENV_GOOGLE_client_sercret); len(a.ClientSecret) == 0 {
		log.Fatalf("%s %s", ENV_GOOGLE_client_sercret, environement_variable_missing)
	}

	// --------- Get RedirectURL
	if a.RedirectURL = os.Getenv(ENV_GOOGLE_redirect_url); len(a.RedirectURL) == 0 {
		log.Fatalf("%s %s", ENV_GOOGLE_redirect_url, environement_variable_missing)
	}

	// --------- Get ScopeEmail
	if a.ScopeEmail = os.Getenv(ENV_GOOGLE_scope_email); len(a.ScopeEmail) == 0 {
		log.Fatalf("%s %s", ENV_GOOGLE_scope_email, environement_variable_missing)
	}

	// --------- Get ScopesUserinfo
	if a.ScopesUserinfo = os.Getenv(ENV_GOOGLE_scope_userinfo); len(a.ScopesUserinfo) == 0 {
		log.Fatalf("%s %s", ENV_GOOGLE_scope_userinfo, environement_variable_missing)
	}

	// --------- Get UserInfoURL
	if a.UserInfoURL = os.Getenv(ENV_GOOGLE_userinfo_url); len(a.UserInfoURL) == 0 {
		log.Fatalf("%s %s", ENV_GOOGLE_userinfo_url, environement_variable_missing)
	}
}
