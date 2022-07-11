package external

import (
	"net/http"
)

func NewHttpServer(http_addr string) error {
	http.HandleFunc("/redirect_42", redirect42)
	http.HandleFunc("/login_google", loginGoogle)
	http.HandleFunc("/callback_google", callbackGoogle)

	if err := http.ListenAndServe(http_addr, nil); err != nil {
		return err
	}

	return nil
}
