package server

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/trixky/hypertube/api-scrapper/proto"
	st "github.com/trixky/hypertube/api-scrapper/sites"
	grpcMetadata "google.golang.org/grpc/metadata"
)

var categories []string = []string{"movies", "shows"}

func (s *ScrapperServer) ScrapeAll(request *pb.ScrapeRequest, out pb.ScrapperService_ScrapeAllServer) error {
	log.Printf("Scrap All %v\n", request)

	for _, scrapper := range st.Scrappers {
		for _, category := range categories {
			var page uint32 = 1
			for {
				page_result, err := scrapper.ScrapeList(category, page)
				if err == nil {
					if err := out.Send(&pb.ScrapeResponse{
						MsDuration: 0,
						Torrents:   page_result.Torrents,
					}); err != nil {
						return err
					}
				} else {
					return err
				}
				page = page_result.NextPage
				if page == 0 {
					break
				}
				time.Sleep(time.Second)
			}
		}
	}

	return nil
}

func (s *ScrapperServer) IdentifyAll(request *pb.IdentifyRequest, out pb.ScrapperService_IdentifyAllServer) error {
	log.Printf("Identify All %v\n", request)
	return nil
}

func (s *ScrapperServer) ScrapeLatest(request *pb.ScrapeLatestRequest, out pb.ScrapperService_ScrapeLatestServer) error {
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
