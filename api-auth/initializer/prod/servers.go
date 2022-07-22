package initializer

import (
	"log"
	"strconv"

	"github.com/trixky/hypertube/api-auth/environment"
	"github.com/trixky/hypertube/api-auth/external"
	"github.com/trixky/hypertube/api-auth/internal"
	"google.golang.org/grpc"
)

func InitServers() *grpc.Server {
	// ------------- GRPC
	grpc_port := ":" + strconv.Itoa(environment.E.GrpcPort)
	grpc_addr := host + grpc_port

	log.Printf("start to serve grpc services on \t\t%s\n", grpc_port)
	grpc_server := internal.NewGrpcServer(grpc_addr)

	// ------------- GRPC-GATEWAY
	grpc_gateway_port := ":" + strconv.Itoa(environment.E.GrpcGatewayPort)

	log.Printf("start to serve grpc gateway services on \t%s\n", grpc_gateway_port)
	internal.NewGrpcGatewayServer(grpc_gateway_port, grpc_addr)

	// ------------- HHTP
	http_addr := ":" + strconv.Itoa(environment.E.HttpPort)

	log.Printf("start to serve http services on \t\t%s\n", http_addr)
	external.NewHttpServer(http_addr)

	return grpc_server
}
