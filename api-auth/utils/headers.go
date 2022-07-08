package utils

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type CookieMe struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	External  string `json:"external"`
}

func HeaderCookieTokenGeneration(token string) *http.Cookie {
	return &http.Cookie{
		Name:  "token",
		Value: token,
	}
}

func HeaderCookieMeGeneration(cookie_me CookieMe, base_64 bool) (*http.Cookie, error) {
	json_value, err := json.Marshal(cookie_me)

	if err != nil {
		return nil, err
	}

	value := string(json_value)

	if base_64 {
		value = base64.StdEncoding.EncodeToString([]byte(value))
	}

	return &http.Cookie{
		Name:  "me",
		Value: value,
	}, nil
}
