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

func (a *ApiGoogle) GetAll() {
	// --------- get ClientId
	if a.ClientId = os.Getenv(ENV_GOOGLE_client_id); len(a.ClientId) == 0 {
		log.Fatalf("%s environement variable missing", ENV_GOOGLE_client_id)
	}

	// --------- get ClientSecret
	if a.ClientSecret = os.Getenv(ENV_GOOGLE_client_sercret); len(a.ClientSecret) == 0 {
		log.Fatalf("%s environement variable missing", ENV_GOOGLE_client_sercret)
	}

	// --------- get RedirectURL
	if a.RedirectURL = os.Getenv(ENV_GOOGLE_redirect_url); len(a.RedirectURL) == 0 {
		log.Fatalf("%s environement variable missing", ENV_GOOGLE_redirect_url)
	}

	// --------- get ScopeEmail
	if a.ScopeEmail = os.Getenv(ENV_GOOGLE_scope_email); len(a.ScopeEmail) == 0 {
		log.Fatalf("%s environement variable missing", ENV_GOOGLE_scope_email)
	}

	// --------- get ScopesUserinfo
	if a.ScopesUserinfo = os.Getenv(ENV_GOOGLE_scope_userinfo); len(a.ScopesUserinfo) == 0 {
		log.Fatalf("%s environement variable missing", ENV_GOOGLE_scope_userinfo)
	}

	// --------- get UserInfoURL
	if a.UserInfoURL = os.Getenv(ENV_GOOGLE_userinfo_url); len(a.UserInfoURL) == 0 {
		log.Fatalf("%s environement variable missing", ENV_GOOGLE_userinfo_url)
	}
}
