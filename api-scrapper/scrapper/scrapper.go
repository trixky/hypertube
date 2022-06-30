package scrapper

import (
	pb "github.com/trixky/hypertube/api-scrapper/proto"
)

type Urls struct {
	Movies, Shows string
}

type ScrapperPageResult struct {
	Torrents []*pb.UnprocessedTorrent
	NextPage uint32
}

type Scrapper struct {
	ScrapeList   func(page_type string, page uint32) (ScrapperPageResult, error)
	ScrapeSingle func(id string) (pb.UnprocessedTorrent, error)
}
