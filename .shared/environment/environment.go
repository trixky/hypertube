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
	DEFAULT_HOST = "0.0.0.0"

	environement_variable_missing = "environement variable missing"
)

type Config struct {
	ENV_grpc_port         string
	ENV_grpc_gateway_port string
	ENV_http_port         string
}

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
