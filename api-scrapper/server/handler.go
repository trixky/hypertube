package server

import (
	"context"
	"fmt"
	"log"

	pb "github.com/trixky/hypertube/api-scrapper/proto"
	grpcMetadata "google.golang.org/grpc/metadata"
)

func (s *ScrapperServiceServer) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Printf("Greet function was invoked with %v\n", in)

	// -----------------------

	md, ok := grpcMetadata.FromIncomingContext(ctx)

	if !ok {
		log.Println("missing args")
		return nil, nil
	}

	scrapper := md.Get("scrapper")

	fmt.Println("scrapper:", scrapper)

	return &pb.SearchResponse{
		Jwt:         "ouiiiiiii",
		Id:          1,
		Page:        1,
		Name:        "Test",
		Genres:      "movie,test",
		Description: "It's a test",
		Thumbnail:   "https://thumbna.il/test",
		Duration:    "1h35min",
		Year:        2022,
		Rating:      72,
	}, nil
}
