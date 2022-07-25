package external

import (
	"fmt"
	"log"

	"github.com/trixky/hypertube/.shared/environment"
	pb "github.com/trixky/hypertube/api-media/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ENV_scrapper_grpc_port = "API_SCRAPPER_GRPC_PORT"
)

type env_scrapper struct {
	Port int
}

// GetAll read all needed enviornment variables
func (e *env_scrapper) GetAll() {
	// --------- Get Scrapper Port
	if grpc_port, err := environment.ReadPort(ENV_scrapper_grpc_port); err != nil {
		log.Fatal(err)
	} else {
		e.Port = grpc_port
	}
}

var Scrapper = env_scrapper{}

var ApiScrapper pb.ScrapperServiceClient

func NewApiScrapperClient() (conn *grpc.ClientConn, err error) {
	scrapper_port := fmt.Sprint(Scrapper.Port)
	conn, err = grpc.Dial("api-scrapper:"+scrapper_port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	ApiScrapper = pb.NewScrapperServiceClient(conn)
	return
}
