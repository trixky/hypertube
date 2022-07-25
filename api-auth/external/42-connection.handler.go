package external

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/trixky/hypertube/.shared/databases"
	"github.com/trixky/hypertube/.shared/environment"
	"github.com/trixky/hypertube/.shared/sanitizer"
	"github.com/trixky/hypertube/.shared/utils"
	"github.com/trixky/hypertube/api-auth/queries"
	"github.com/trixky/hypertube/api-auth/sqlc"
)

type code42Response struct {
	Access_token  string `json:"access_token"`
	Token_type    string `json:"Token_type"`
	Expires_in    int    `json:"expires_in"`
	Refresh_token string `json:"refresh_token"`
	Scope         string `json:"scope"`
	Created_at    int    `json:"created_at"`
}

type me42Response struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	Login     string `json:"login"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
}

// sanitizeExternal42Connection sanitizes inputs for "redirect42" route
func sanitizeExternal42Connection(code string) error {
	if err := sanitizer.Sanitize42Code(code); err != nil { // 42 code
		return err
	}

	return nil
}

// getTokenFromCode get 42 API token from the code
func getTokenFromCode(code string) (*code42Response, error) {
	// https://profile.intra.42.fr/oauth/applications/9554
	// https://api.intra.42.fr/apidoc/guides/web_application_flow
	// https://api.intra.42.fr/apidoc/guides/getting_started
	// https://api.intra.42.fr/apidoc

	// Encode the data for the 42 API
	response, err := http.PostForm(environment.Api42.RequestUrl, url.Values{
		"grant_type":    {environment.Api42.GrantType},
		"client_id":     {environment.Api42.ClientId},
		"client_secret": {environment.Api42.ClientSecret},
		"code":          {code},
		"redirect_uri":  {environment.Api42.RedirectionUri},
	})

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	// Read the body of the response
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	code_42_response := code42Response{}

	// Extract JSON from the body
	if err := json.Unmarshal(body, &code_42_response); err != nil {
		return nil, err
	}

	return &code_42_response, nil
}

// getTokenFromCode get the 42 user information from his token
func getUser42InfoFromToken(code_42_response *code42Response) (*me42Response, error) {
	req, err := http.NewRequest("GET", environment.Api42.RequestMe, nil)

	if err != nil {
		return nil, err
	}

	// Generates the authorization header from the 42 user token
	authorization_header := code_42_response.Token_type + " " + code_42_response.Access_token

	req.Header.Add("Authorization", authorization_header)

	client := &http.Client{}

	// Send th request
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Read the body of the response
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	me_42_response := me42Response{}

	// Extract JSON from the body
	if err := json.Unmarshal(body, &me_42_response); err != nil {
		return nil, err
	}

	return &me_42_response, nil
}

// Redirect42 Handles the "redirect42" route
func redirect42(w http.ResponseWriter, r *http.Request) {
	// -------------------- Sanitize
	code := r.URL.Query().Get("code")
	if err := sanitizeExternal42Connection(code); err != nil {
		http.Error(w, "code missing", http.StatusBadRequest)
		return
	}

	// -------------------- 42 api call
	code_42_response, err := getTokenFromCode(code)
	if err != nil {
		http.Error(w, "token extration failed", http.StatusInternalServerError)
		return
	}

	me_42_response, err := getUser42InfoFromToken(code_42_response)
	if err != nil {
		http.Error(w, "user info extraction failed", http.StatusInternalServerError)
		return
	}

	// -------------------- DB
	user, err := queries.SqlcQueries.Create42ExternalUser(context.Background(), sqlc.Create42ExternalUserParams{
		Email:     me_42_response.Email,
		Username:  me_42_response.Login,
		Firstname: me_42_response.Firstname,
		Lastname:  me_42_response.Lastname,
		ID42: sql.NullInt32{
			Int32: int32(me_42_response.Id),
			Valid: true,
		},
	})

	if err != nil {
		if databases.ErrorIsDuplication(err) {
			user, err = queries.SqlcQueries.GetUserBy42Id(context.Background(), sql.NullInt32{
				Int32: int32(me_42_response.Id),
				Valid: true,
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
	// Generates the token
	token := uuid.New().String()

	if err := queries.AddToken(user.ID, token, databases.REDIS_EXTERNAL_42); err != nil {
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
		External:  databases.REDIS_EXTERNAL_42,
	}, true)

	if err != nil {
		http.Error(w, "cookie generation failed", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, me)

	http.Redirect(w, r, "http://"+environment.Client.Domain+":"+fmt.Sprint(environment.Client.Port)+"/", http.StatusFound)
}
