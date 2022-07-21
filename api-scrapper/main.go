package main

import (
	"log"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/trixky/hypertube/api-scrapper/databases"
	"github.com/trixky/hypertube/api-scrapper/environment"
	"github.com/trixky/hypertube/api-scrapper/internal"
)

const (
	host            = "0.0.0.0"
	postgres_driver = "postgres"
)

func main() {
	environment.E.GetAll()

	// ------------- postgres
	log.Printf("start connection to postgres on %s:%d\n", environment.E.PostgresHost, environment.E.PostgresPort)
	if err := databases.InitPosgres(databases.PostgresConfig{
		Driver:   postgres_driver,
		Host:     environment.E.PostgresHost,
		Port:     environment.E.PostgresPort,
		User:     environment.E.PostgresUser,
		Password: environment.E.PostgresPassword,
		Dbname:   environment.E.PostgresDB,
	}); err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	// ------------- grpc
	grpc_addr := host + ":" + strconv.Itoa(environment.E.GrpcPort)
	go func() {
		log.Fatalf("failed to serve grpc on: %v\n", internal.NewGrpcServer(grpc_addr))
	}()

	// ------------- loop forever to scrape
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.SingletonMode()
	scheduler.StartImmediately()
	scheduler.Every(30).Minutes().Do(func() {
		internal.DoScrapeLatest(nil)
		log.Println("Next scrape in 30min")
	})
	scheduler.StartBlocking()
}
