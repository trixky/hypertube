package external

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/trixky/hypertube/api-auth/databases"
	"github.com/trixky/hypertube/api-auth/environment"
	"github.com/trixky/hypertube/api-auth/sanitizer"
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

func sanitizeExternal42Connection(code string) error {
	if err := sanitizer.Sanitize42Code(code); err != nil { // 42 code
		return err
	}

	return nil
}

func getTokenFromCode(code string) (*code42Response, error) {
	// https://profile.intra.42.fr/oauth/applications/9554
	// https://api.intra.42.fr/apidoc/guides/web_application_flow
	// https://api.intra.42.fr/apidoc/guides/getting_started
	// https://api.intra.42.fr/apidoc

	//Encode the data
	response, err := http.PostForm(environment.E.API42.RequestUrl, url.Values{
		"grant_type":    {environment.E.API42.GrantType},
		"client_id":     {environment.E.API42.ClientId},
		"client_secret": {environment.E.API42.ClientSecret},
		"code":          {code},
		"redirect_uri":  {environment.E.API42.RedirectionUri},
	})

	//okay, moving on...
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	code_42_response := code42Response{}

	if err := json.Unmarshal(body, &code_42_response); err != nil {
		return nil, err
	}

	return &code_42_response, nil
}

func getUser42InfoFromToken(code_42_response *code42Response) (*me42Response, error) {
	req, err := http.NewRequest("GET", environment.E.API42.RequestMe, nil)

	if err != nil {
		return nil, err
	}

	authorization_header := code_42_response.Token_type + " " + code_42_response.Access_token

	req.Header.Add("Authorization", authorization_header)

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	me_42_response := me42Response{}

	if err := json.Unmarshal(body, &me_42_response); err != nil {
		return nil, err
	}

	return &me_42_response, nil
}

func redirect42(w http.ResponseWriter, r *http.Request) {
	// -------------------- sanitize
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

	// -------------------- db
	user, err := databases.DBs.SqlcQueries.Create42ExternalUser(context.Background(), sqlc.Create42ExternalUserParams{
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
			user, err = databases.DBs.SqlcQueries.GetUserBy42Id(context.Background(), sql.NullInt32{
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

	// -------------------- cache
	token := uuid.New().String() // token generation

	if err := databases.AddToken(user.ID, token, databases.REDIS_EXTERNAL_42); err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:  "token",
		Value: token,
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "http://localhost:4040/", http.StatusFound)
}
