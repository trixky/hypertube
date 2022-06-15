package main

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/trixky/hypertube/api-scrapper/postgres"
	pb "github.com/trixky/hypertube/api-scrapper/proto"
	"github.com/trixky/hypertube/api-scrapper/server"
	"google.golang.org/grpc"
)

const (
	host            = "0.0.0.0"
	database_driver = "postgres"

	env_port              = "API_AUTH_PORT"
	env_postgres_host     = "POSTGRES_HOST"
	env_postgres_port     = "POSTGRES_PORT"
	env_postgres_user     = "POSTGRES_USER"
	env_postgres_password = "POSTGRES_PASSWORD"
	env_postgres_db       = "POSTGRES_DB"
)

type Env struct {
	Port             int
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
}

func (e *Env) GetAll() {
	const (
		port_max = 65535
		port_min = 1000
	)

	// --------- get Port
	port, err := strconv.Atoi(os.Getenv(env_port))

	if err != nil {
		log.Fatalf("%s environement variable is corrupted or missing", env_port)
	} else if port < port_min || port > port_max {
		log.Fatalf("port need to be included between %d and %d (%d)", port_min, port_max, port)
	}

	e.Port = port

	// --------- get PostgresHost
	if e.PostgresHost = os.Getenv(env_postgres_host); len(e.PostgresHost) == 0 {
		log.Fatalf("%s environement variable missing", env_postgres_host)
	}

	// --------- get PostgresPort
	postgres_port, err := strconv.Atoi(os.Getenv(env_postgres_port))

	if err != nil {
		log.Fatalf("%s environement variable is corrupted or missing", env_postgres_port)
	} else if postgres_port < port_min || postgres_port > port_max {
		log.Fatalf("port need to be included between %d and %d (%d)", port_min, port_max, postgres_port)
	}

	e.PostgresPort = postgres_port

	// --------- get PostgresUser
	if e.PostgresUser = os.Getenv(env_postgres_user); len(e.PostgresUser) == 0 {
		log.Fatalf("%s environement variable missing", env_postgres_user)
	}

	// --------- get PostgresPassword
	if e.PostgresPassword = os.Getenv(env_postgres_password); len(e.PostgresPassword) == 0 {
		log.Fatalf("%s environement variable missing", env_postgres_password)
	}

	// --------- get PostgresDB
	if e.PostgresDB = os.Getenv(env_postgres_db); len(e.PostgresDB) == 0 {
		log.Fatalf("%s environement variable missing", env_postgres_db)
	}
}

func readEnv() (env Env) {
	env.GetAll()

	return
}

func main() {
	env := readEnv()

	addr := host + ":" + strconv.Itoa(env.Port)

	if err := postgres.Init(postgres.Config{
		Driver:   database_driver,
		Host:     env.PostgresHost,
		Port:     env.PostgresPort,
		User:     env.PostgresUser,
		Password: env.PostgresPassword,
		Dbname:   env.PostgresDB,
	}); err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	log.Println("connected to postgres")

	listen, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("failed to listen on: %v\n", err)
	}

	log.Printf("listening on %s\n", addr)

	s := grpc.NewServer()

	pb.RegisterAuthServiceServer(s, &server.AuthServer{})

	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve on: %v\n", err)
	}
}
