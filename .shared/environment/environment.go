package environment

// https://github.com/grpc-ecosystem/grpc-gateway/blob/f046a4ebdc9be76e11c6239eaeba4e30e9e2444e/docs/docs/faq.md
// https://github.com/grpc-ecosystem/grpc-gateway/blob/master/examples/internal/gateway/main.go
// https://github.com/grpc-ecosystem/grpc-gateway/blob/0ebdfba80649c56b0da0777376c970d17c3c9540/examples/internal/gateway/handlers.go#L32

import (
	"fmt"
	"os"
	"strconv"
)

const (
	environement_variable_missing = "environement variable missing"
)

// type Env struct {
// 	GRPC          Grpc
// 	HTTP          Http
// 	PG            Postgres
// 	RD            Redis
// 	API42         Api42
// 	APIGoogle     ApiGoogle
// 	OUTLOOKConfig OutlookConfig
// }

// var ENV = Env{}

// read_port convert and sanitize string port to integer
func read_port(name string) (int, error) {
	const (
		port_max = 65535
		port_min = 1000
	)

	port, err := strconv.Atoi(os.Getenv(name))

	if err != nil {
		return 0, fmt.Errorf("%s environement variable is corrupted or missing", name)
	} else if port < port_min || port > port_max {
		return 0, fmt.Errorf("port need to be included between %d and %d (%d)", port_min, port_max, port)
	}

	return port, nil
}

func ReadAll() {
	Api42.GetAll()
	ApiGoogle.GetAll()
	Grpc.GetAll()
	Http.GetAll()
	Outlook.GetAll()
	Postgres.GetAll()
	Redis.GetAll()
}

// func (e *Env) GetAll() {
// 	e.GRPC.GetAll()          // Get http environment variables
// 	e.HTTP.GetAll()          // Get http environment variables
// 	e.PG.GetAll()            // Get postgres environment variables
// 	e.RD.GetAll()            // Get redis environment variables
// 	e.API42.GetAll()         // Get 42 environment variables
// 	e.APIGoogle.GetAll()     // Get google environment variables
// 	e.OUTLOOKConfig.GetAll() // Get outlook environment variables
// }
