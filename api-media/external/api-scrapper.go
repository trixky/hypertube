package external

import (
	"fmt"

	"github.com/trixky/hypertube/api-media/environment"
	pb "github.com/trixky/hypertube/api-media/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ApiScrapper pb.ScrapperServiceClient

func NewApiScrapperClient() (conn *grpc.ClientConn, err error) {
	scrapper_port := fmt.Sprint(environment.E.ScrapperGrpcPort)
	conn, err = grpc.Dial("api-scrapper:"+scrapper_port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	ApiScrapper = pb.NewScrapperServiceClient(conn)
	return
}
