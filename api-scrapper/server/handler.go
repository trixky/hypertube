package server

import (
	"context"
	"fmt"
	"log"

	"github.com/gocolly/colly"
	pb "github.com/trixky/hypertube/api-scrapper/proto"
	grpcMetadata "google.golang.org/grpc/metadata"
)

func (s *ScrapperServer) ScrapAll(request *pb.ScrapRequest, out pb.ScrapperService_ScrapAllServer) error {
	log.Printf("Scrap All %v\n", request)

	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	// Find and print all links
	c.OnHTML(".mw-parser-output", func(e *colly.HTMLElement) {
		links := e.ChildAttrs("a", "href")
		fmt.Println(links)
	})
	c.Visit("https://en.wikipedia.org/wiki/Web_scraping")

	var i uint32 = 0
	for ; i < 5; i++ {
		if err := out.Send(&pb.ScrapResponse{
			MsDuration: i,
			Torrents:   []*pb.Torrent{},
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *ScrapperServer) IdentifyAll(request *pb.IdentifyRequest, out pb.ScrapperService_IdentifyAllServer) error {
	log.Printf("Identify All %v\n", request)
	return nil
}

func (s *ScrapperServer) ScrapLatest(request *pb.ScrapLatestRequest, out pb.ScrapperService_ScrapLatestServer) error {
	log.Printf("Scrap Latest %v\n", request)
	return nil
}

func (s *ScrapperServer) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	md, ok := grpcMetadata.FromIncomingContext(ctx)

	if !ok {
		log.Println("missing args")
		return nil, nil
	}

	search := md.Get("search")
	fmt.Println("search:", search)

	return &pb.SearchResponse{
		Page:   1,
		Medias: []*pb.Media{},
	}, nil
}

func (s *ScrapperServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	md, ok := grpcMetadata.FromIncomingContext(ctx)

	if !ok {
		log.Println("missing args")
		return nil, nil
	}

	search := md.Get("get")
	fmt.Println("get:", search)

	return &pb.GetResponse{
		Id:          1,
		ImdbId:      "tt3456",
		Name:        "Movie",
		Description: "Movie ?",
		Year:        2000,
		TorrentPublicInformations: &pb.TorrentPublicInformations{
			Name:  "Movie [1080p]",
			Seed:  &[]uint32{42}[0],
			Leech: &[]uint32{42}[0],
			Size:  &[]string{"123456789"}[0],
		},
		Staffs: []*pb.Staff{
			&pb.Staff{
				Id:        1,
				ImdbId:    "tt1234",
				Name:      "Writer",
				Role:      "Writer",
				Thumbnail: "jpg",
				Url:       "http",
			},
		},
		Relations: []*pb.Relation{
			&pb.Relation{
				Id:        2,
				ImdbId:    "tt2345",
				Name:      "Movie 2",
				Thumbnail: "jpg",
			},
		},
		Duration:   &[]string{"1h42"}[0],
		Thumbnail:  &[]string{"jpg"}[0],
		Background: &[]string{"jpg"}[0],
		Genres:     &[]string{"movie,action"}[0],
		Rating:     &[]string{"80"}[0],
	}, nil
}
