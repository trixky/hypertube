package initializer

import (
	"log"
	"strconv"

	"github.com/trixky/hypertube/.shared/environment"
	"github.com/trixky/hypertube/api-auth/external"
	"github.com/trixky/hypertube/api-auth/internal"
	"google.golang.org/grpc"
)

const host = "0.0.0.0"

func InitServers() *grpc.Server {
	// ------------- GRPC
	grpc_port := ":" + strconv.Itoa(environment.Grpc.Port)
	grpc_addr := host + grpc_port

	log.Printf("start to serve grpc services on \t\t%s\n", grpc_port)
	grpc_server := internal.NewGrpcServer(grpc_addr)

	// ------------- GRPC-GATEWAY
	grpc_gateway_port := ":" + strconv.Itoa(environment.Grpc.GatewayPort)

	log.Printf("start to serve grpc gateway services on \t%s\n", grpc_gateway_port)
	internal.NewGrpcGatewayServer(grpc_gateway_port, grpc_addr)

	// ------------- HHTP
	http_addr := ":" + strconv.Itoa(environment.Http.Port)

	log.Printf("start to serve http services on \t\t%s\n", http_addr)
	external.NewHttpServer(http_addr)

	return grpc_server
}
