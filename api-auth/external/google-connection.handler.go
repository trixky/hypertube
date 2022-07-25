// https://www.youtube.com/watch?v=Vmi3trk0rCk

package external

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/trixky/hypertube/.shared/databases"
	"github.com/trixky/hypertube/.shared/environment"
	"github.com/trixky/hypertube/.shared/utils"
	"github.com/trixky/hypertube/api-auth/queries"
	"github.com/trixky/hypertube/api-auth/sqlc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleConfig *oauth2.Config // = setupGoogleConfig()
var googleLoginUrl string       // = generateGoogleLoginUrl()

const (
	QUERY_KEY_state           = "state"
	QUERY_KEY_code            = "code"
	QUERY_VALUE_state_default = "default"
)

type meGoogleResponse struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	Login     string `json:"name"`
	Firstname string `json:"given_name"`
	Lastname  string `json:"family_name"`
}

// setupGoogleConfig setups the google config from environment variables
func setupGoogleConfig() {
	config := &oauth2.Config{
		ClientID:     environment.ApiGoogle.ClientId,
		ClientSecret: environment.ApiGoogle.ClientSecret,
		RedirectURL:  environment.ApiGoogle.RedirectURL,
		Scopes: []string{
			environment.ApiGoogle.ScopeEmail,
			environment.ApiGoogle.ScopesUserinfo,
		},
		Endpoint: google.Endpoint,
	}

	googleConfig = config
}

// generateGoogleLoginUrl generates the google login url from the config
func generateGoogleLoginUrl() {
	if googleConfig == nil {
		setupGoogleConfig()
	}
	googleLoginUrl = googleConfig.AuthCodeURL(QUERY_VALUE_state_default)
}

// loginGoogle logs in to google from the google login url
func loginGoogle(w http.ResponseWriter, r *http.Request) {
	if googleLoginUrl == "" {
		generateGoogleLoginUrl()
	}

	http.Redirect(w, r, googleLoginUrl, http.StatusSeeOther)
}

// callbackGoogle Handles the "callbackGoogle" route
func callbackGoogle(w http.ResponseWriter, r *http.Request) {
	query_values := r.URL.Query()

	state := query_values.Get(QUERY_KEY_state)

	if len(state) == 0 {
		http.Error(w, QUERY_KEY_state+" missing", http.StatusBadRequest)
		return
	}

	// Get the code from the url query
	code := query_values.Get(QUERY_KEY_code)

	if len(code) == 0 {
		http.Error(w, QUERY_KEY_code+" missing", http.StatusBadRequest)
		return
	}

	// Generates the google login url
	if googleConfig == nil {
		generateGoogleLoginUrl()
	}

	// Exchange the code to the google token
	google_token, err := googleConfig.Exchange(context.Background(), code)

	if err != nil {
		http.Error(w, QUERY_KEY_code+" exchange failed", http.StatusForbidden)
		return
	}

	// Exchange the google token to the google user info
	google_r, err := http.Get(environment.ApiGoogle.UserInfoURL + "?access_token=" + google_token.AccessToken)

	if err != nil {
		http.Error(w, "user info scrapping failed", http.StatusInternalServerError)
		return
	}

	// Read the body of the response
	user_data, err := ioutil.ReadAll(google_r.Body)

	if err != nil {
		http.Error(w, "user info reading failed 1", http.StatusInternalServerError)
		return
	}

	me_google_response := meGoogleResponse{}

	// Extract JSON from the body
	if err := json.Unmarshal(user_data, &me_google_response); err != nil {
		http.Error(w, "user info reading failed 2", http.StatusInternalServerError)
	}

	// -------------------- DB
	user, err := queries.SqlcQueries.CreateGoogleExternalUser(context.Background(), sqlc.CreateGoogleExternalUserParams{
		Email:     me_google_response.Email,
		Username:  me_google_response.Login,
		Firstname: me_google_response.Firstname,
		Lastname:  me_google_response.Lastname,
		IDGoogle: sql.NullString{
			String: me_google_response.Id,
			Valid:  true,
		},
	})

	if err != nil {
		if databases.ErrorIsDuplication(err) {
			user, err = queries.SqlcQueries.GetUserByGoogleId(context.Background(), sql.NullString{
				String: me_google_response.Id,
				Valid:  true,
			})

			if err != nil {
				http.Error(w, "user retrieve failed", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "user creation failed", http.StatusInternalServerError)
			return
		}
	}

	// -------------------- Cache
	token := uuid.New().String() // token generation

	if err := queries.AddToken(user.ID, token, databases.REDIS_EXTERNAL_google); err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	token_cookie := utils.HeaderCookieTokenGeneration(token)

	http.SetCookie(w, token_cookie)

	me, err := utils.HeaderCookieUserGeneration(utils.User{
		Id:        int(user.ID),
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		External:  databases.REDIS_EXTERNAL_google,
	}, true)

	if err != nil {
		http.Error(w, "cookie generation failed", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, me)

	http.Redirect(w, r, "http://"+environment.Client.Domain+":"+fmt.Sprint(environment.Client.Port)+"/", http.StatusFound)
}
