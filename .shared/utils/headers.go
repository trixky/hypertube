package utils

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	External  string `json:"external"`
}

const (
	COOKIE_KEY_user_info = "userInfo"
	COOKIE_KEY_token     = "token"
)

// HeaderCookieTokenGeneration generates/insert the token in cookie header
func HeaderCookieTokenGeneration(token string) *http.Cookie {
	return &http.Cookie{
		Name:  COOKIE_KEY_token,
		Value: token,
	}
}

// HeaderCookieUserGeneration generates/insert the user infos in cookie header
func HeaderCookieUserGeneration(cookie_user User, base_64 bool) (*http.Cookie, error) {
	// JSON the user infos
	json_value, err := json.Marshal(cookie_user)

	if err != nil {
		return nil, err
	}

	value := string(json_value)

	// Convert JSON to BASE 64
	if base_64 {
		value = base64.StdEncoding.EncodeToString([]byte(value))
	}

	return &http.Cookie{
		Name:  COOKIE_KEY_user_info,
		Value: value,
	}, nil
}
