package server

import (
	"context"
	"fmt"
	"log"

	"github.com/gocolly/colly"
	pb "github.com/trixky/hypertube/api-scrapper/proto"
	grpcMetadata "google.golang.org/grpc/metadata"
)

func (s *ScrapperServer) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Printf("Greet function was invoked with %v\n", in)

	md, ok := grpcMetadata.FromIncomingContext(ctx)

	if !ok {
		log.Println("missing args")
		return nil, nil
	}

	search := md.Get("search")
	fmt.Println("search:", search)

	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	// Find and print all links
	c.OnHTML(".mw-parser-output", func(e *colly.HTMLElement) {
		links := e.ChildAttrs("a", "href")
		fmt.Println(links)
	})
	c.Visit("https://en.wikipedia.org/wiki/Web_scraping")

	return &pb.SearchResponse{
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
