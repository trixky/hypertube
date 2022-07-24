package main

import (
	"log"

	"github.com/trixky/hypertube/.shared/environment"
	"github.com/trixky/hypertube/api-auth/databases"
	"github.com/trixky/hypertube/api-auth/initializer"
)

func init() {
	log.Println("------------------------- INIT api-auth")
	environment.ReadAll()     // Get environment variables
	databases.InitDBs()       // Init DBs
	initializer.InitServers() // Init servers
}

func main() {
	log.Println("------------------------- START api-auth")
	select {} // Keep alive
}
