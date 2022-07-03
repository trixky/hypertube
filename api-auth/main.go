package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/trixky/hypertube/api-auth/databases"
	pb "github.com/trixky/hypertube/api-auth/proto"
	"github.com/trixky/hypertube/api-auth/server"
	"github.com/trixky/hypertube/api-auth/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	host            = "0.0.0.0"
	postgres_driver = "postgres"
)

func main() {
	env := utils.ReadEnv()

	grpc_addr := host + ":" + strconv.Itoa(env.GrpcPort)

	if err := databases.InitPosgres(databases.PostgresConfig{
		Driver:   postgres_driver,
		Host:     env.PostgresHost,
		Port:     env.PostgresPort,
		User:     env.PostgresUser,
		Password: env.PostgresPassword,
		Dbname:   env.PostgresDB,
	}); err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	log.Println("connected to postgres")

	if err := databases.InitRedis(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	log.Println("connected to redis")

	listen, err := net.Listen("tcp", grpc_addr)

	if err != nil {
		log.Fatalf("failed to listen on: %v\n", err)
	}

	log.Printf("listening on %s\n", grpc_addr)

	s := grpc.NewServer()

	pb.RegisterAuthServiceServer(s, &server.AuthServer{})

	go func() {
		log.Fatalf("failed to serve on: %v\n", s.Serve(listen))
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		grpc_addr,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	http_addr := ":" + strconv.Itoa(env.HttpPort)

	gwmux := runtime.NewServeMux(runtime.WithMetadata(
		func(ctx context.Context, r *http.Request) metadata.MD {
			log.Println("authauth:", r.Header.Get("authauth"))

			md := make(map[string]string)
			if method, ok := runtime.RPCMethod(ctx); ok {
				md["method"] = method // /grpc.gateway.examples.internal.proto.examplepb.LoginService/Login
			}
			if pattern, ok := runtime.HTTPPathPattern(ctx); ok {
				md["pattern"] = pattern // /v1/example/login
			}
			return metadata.New(md)
		},
	))

	err = pb.RegisterAuthServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    http_addr,
		Handler: utils.AllowCORS(gwmux),
	}

	// -------------
	// mux := http.NewServeMux()
	// mux.Handle("/", gwmux)
	// -------------

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090 /", http_addr)
	log.Fatalln(gwServer.ListenAndServe())

}
