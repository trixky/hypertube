package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/trixky/hypertube/.shared/environment"
	"github.com/trixky/hypertube/api-scrapper/databases"
	"github.com/trixky/hypertube/api-scrapper/internal"
)

func init() {
	log.Println("------------------------- INIT api-scrapper")

	// Set environment config
	environment_config := environment.Config{
		ENV_grpc_port: "API_SCRAPPER_GRPC_PORT",
	}

	environment.Postgres.GetAll()                // Get postgres environment
	environment.Grpc.GetAll(&environment_config) // Get grpc environment
	environment.TMDB.GetAll()                    // Get TMDB API key

	databases.InitDBs()      // Init DBs
	internal.NewGrpcServer() // Init internal servers
}

func main() {
	log.Println("------------------------- START api-scrapper")
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
