package main

import (
	"log"

	"github.com/trixky/hypertube/api-auth/databases"
	"github.com/trixky/hypertube/api-auth/environment"
	initializer "github.com/trixky/hypertube/api-auth/initializer"
)

func init() {
	log.Println("------------------------- INIT api-auth")
	environment.E.GetAll()    // Get environment variables
	databases.InitDBs()       // Init DBs
	initializer.InitServers() // Init servers
}

func main() {
	log.Println("------------------------- START api-auth")
	select {} // Keep alive
}
