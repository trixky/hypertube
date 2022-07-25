package external

import (
	"log"
	"net/http"
	"strconv"

	"github.com/trixky/hypertube/.shared/environment"
)

// NewHttpServer create a new HTTP server
func NewHttpServer() {
	http_addr := ":" + strconv.Itoa(environment.Http.Port)

	log.Printf("start to serve http services on \t\t%s\n", http_addr)

	http.HandleFunc("/redirect_42", redirect42)
	http.HandleFunc("/login_google", loginGoogle)
	http.HandleFunc("/callback_google", callbackGoogle)

	go func() {
		log.Fatalf("failed to serve grpc-gateway: %v\n", http.ListenAndServe(http_addr, nil))
	}()
}
