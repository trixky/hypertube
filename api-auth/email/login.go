// https://gist.github.com/andelf/5118732

package email

import (
	"errors"
	"net/smtp"
)

type auth_outlook struct {
	username, password string
}

func loginAuthOutlook(username, password string) smtp.Auth {
	return &auth_outlook{username, password}
}

func (a *auth_outlook) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *auth_outlook) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}
