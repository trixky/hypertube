package internal

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	pb "github.com/trixky/hypertube/api-scrapper/proto"
	st "github.com/trixky/hypertube/api-scrapper/sites"
)

func (s *ScrapperServer) ScrapeAll(request *pb.ScrapeRequest, out pb.ScrapperService_ScrapeAllServer) error {
	ctx := context.Background()
	log.Println("Scrape All", request)

	for _, scrapper := range st.Scrappers {
		var err error
		for _, category := range Categories {
			var page uint32 = 1
			for {
				start := time.Now()
				page_result, err_scrapper := scrapper.ScrapeList(category, page)
				if err_scrapper != nil {
					err = err_scrapper
					break
				}
				new_torrents := make([]*pb.Torrent, 0, 30)

				// Update each existing torrents or add them to the database
				for _, torrent := range page_result.Torrents {
					created, created_torrent, err := UpdateOrCreateTorrent(ctx, &scrapper, torrent)
					if err != nil {
						return err
					}
					if created {
						converted_torrent := TorrentToProto(&created_torrent)
						new_torrents = append(new_torrents, &converted_torrent)
					}
				}

				// Send the new torrents on the stream
				if err := out.Send(&pb.ScrapeResponse{
					MsDuration: uint32(time.Since(start).Milliseconds()),
					Torrents:   new_torrents,
				}); err != nil {
					return err
				}

				// Update NextPage to loop or complete the job
				page = page_result.NextPage
				if page == 0 {
					break
				}
				time.Sleep(time.Second)
			}
			if err != nil {
				break
			}
		}

		// Handle timeout errors and skip to the next site
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || os.IsTimeout(err) {
				log.Println("Skipped scrapper", scrapper, "on timeout error !")
				err = nil
				continue
			} else {
				return err
			}
		}
	}

	log.Println("Done ScrapeAll")

	return nil
}
