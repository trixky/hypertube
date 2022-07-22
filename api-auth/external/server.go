package external

import (
	"log"
	"net/http"
)

// NewHttpServer create a new HTTP server
func NewHttpServer(http_addr string) {
	http.HandleFunc("/redirect_42", redirect42)
	http.HandleFunc("/login_google", loginGoogle)
	http.HandleFunc("/callback_google", callbackGoogle)

	go func() {
		log.Fatalf("failed to serve grpc-gateway: %v\n", http.ListenAndServe(http_addr, nil))
	}()
}
