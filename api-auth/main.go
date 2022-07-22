package main

import (
	"log"

	"github.com/trixky/hypertube/api-auth/environment"
	initializer "github.com/trixky/hypertube/api-auth/initializer/prod"
)

func init() {
	log.Println("------------------------- INIT api-auth")
	environment.E.GetAll()    // Get environment variables
	initializer.InitDBs()     // Init DBs
	initializer.InitServers() // Init servers
}

func main() {
	log.Println("------------------------- START api-auth")
	select {} // Keep alive
}
