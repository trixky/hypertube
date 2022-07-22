package initializer

import (
	"log"

	databases "github.com/trixky/hypertube/api-auth/databases/mock"
)

const (
	host            = "0.0.0.0"
	postgres_driver = "postgres"
)

func InitDBs() {
	// ------------- POSTGRES
	log.Println("start connection to postgres (mock)")
	databases.InitPostgresMock()

	// ------------- REDIS
	log.Println("start connection to redis (mock)")
	databases.InitRedisMock()
}
