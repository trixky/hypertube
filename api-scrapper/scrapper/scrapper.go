package scrapper

import (
	"regexp"
	"strings"

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
	ScrapeSingle func(torrent *pb.UnprocessedTorrent) error
}

var IMDBre = regexp.MustCompile("(?:imdb\\.com\\/title\\/)?(tt\\d+)")

func ExtractPath(filename string) (*string, string) {
	path := strings.Split(filename, "/")
	var joined_path *string
	if len(path) > 1 {
		last := len(path) - 1
		filename = path[last]
		file_path := strings.Join(path[:last], "/")
		joined_path = &file_path
	}
	return joined_path, filename
}
