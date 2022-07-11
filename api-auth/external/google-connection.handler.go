// https://www.youtube.com/watch?v=Vmi3trk0rCk

package external

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/trixky/hypertube/api-auth/databases"
	"github.com/trixky/hypertube/api-auth/environment"
	"github.com/trixky/hypertube/api-auth/sqlc"
	"github.com/trixky/hypertube/api-auth/utils"
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

func setupGoogleConfig() {
	config := &oauth2.Config{
		ClientID:     environment.E.APIGoogle.ClientId,
		ClientSecret: environment.E.APIGoogle.ClientSecret,
		RedirectURL:  environment.E.APIGoogle.RedirectURL,
		Scopes: []string{
			environment.E.APIGoogle.ScopeEmail,
			environment.E.APIGoogle.ScopesUserinfo,
		},
		Endpoint: google.Endpoint,
	}

	googleConfig = config
}

func generateGoogleLoginUrl() {
	if googleConfig == nil {
		setupGoogleConfig()
	}
	googleLoginUrl = googleConfig.AuthCodeURL(QUERY_VALUE_state_default)
}

func loginGoogle(w http.ResponseWriter, r *http.Request) {
	if googleLoginUrl == "" {
		generateGoogleLoginUrl()
	}

	http.Redirect(w, r, googleLoginUrl, http.StatusSeeOther)
}

func callbackGoogle(w http.ResponseWriter, r *http.Request) {
	query_values := r.URL.Query()

	state := query_values.Get(QUERY_KEY_state)

	if len(state) == 0 {
		http.Error(w, QUERY_KEY_state+" missing", http.StatusBadRequest)
		return
	}

	code := query_values.Get(QUERY_KEY_code)

	if len(code) == 0 {
		http.Error(w, QUERY_KEY_code+" missing", http.StatusBadRequest)
		return
	}

	if googleConfig == nil {
		generateGoogleLoginUrl()
	}

	google_token, err := googleConfig.Exchange(context.Background(), code)

	if err != nil {
		http.Error(w, QUERY_KEY_code+" exchange failed", http.StatusForbidden)
		return
	}
	google_r, err := http.Get(environment.E.APIGoogle.UserInfoURL + "?access_token=" + google_token.AccessToken)

	if err != nil {
		http.Error(w, "user info scrapping failed", http.StatusInternalServerError)
		return
	}

	user_data, err := ioutil.ReadAll(google_r.Body)

	if err != nil {
		http.Error(w, "user info reading failed 1", http.StatusInternalServerError)
		return
	}

	me_google_response := meGoogleResponse{}

	if err := json.Unmarshal(user_data, &me_google_response); err != nil {
		http.Error(w, "user info reading failed 2", http.StatusInternalServerError)
	}

	// -------------------- db
	user, err := databases.DBs.SqlcQueries.CreateGoogleExternalUser(context.Background(), sqlc.CreateGoogleExternalUserParams{
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
			user, err = databases.DBs.SqlcQueries.GetUserByGoogleId(context.Background(), sql.NullString{
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

	// -------------------- cache
	token := uuid.New().String() // token generation

	if err := databases.AddToken(user.ID, token, databases.EXTERNAL_GOOGLE); err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	token_cookie := utils.HeaderCookieTokenGeneration(token)

	http.SetCookie(w, token_cookie)

	me, err := utils.HeaderCookieMeGeneration(utils.User{
		Id:        int(user.ID),
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		External:  databases.EXTERNAL_GOOGLE,
	}, true)

	if err != nil {
		http.Error(w, "cookie generation failed", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, me)

	http.Redirect(w, r, "http://localhost:4040/", http.StatusFound)
}
