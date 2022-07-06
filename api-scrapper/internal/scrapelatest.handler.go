package internal

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	pb "github.com/trixky/hypertube/api-scrapper/proto"
	st "github.com/trixky/hypertube/api-scrapper/sites"
)

func DoScrapeLatest(callback *func(response *pb.ScrapeResponse) error) error {
	ctx := context.Background()
	log.Println("Scrape Latest")

	for _, scrapper := range st.Scrappers {
		var err error
		for _, category := range Categories {
			var page uint32 = 1
			has_existing := false
			for {
				start := time.Now()
				page_result, err_scrapper := scrapper.ScrapeList(category, page)
				if err_scrapper != nil {
					err = err_scrapper
					break
				}

				// Update each existing torrents or add them to the database
				new_torrents := make([]*pb.Torrent, 0, 30)
				for _, torrent := range page_result.Torrents {
					created, created_torrent, err := UpdateOrCreateTorrent(ctx, &scrapper, torrent)
					if err != nil {
						return err
					}
					if created {
						converted_torrent := TorrentToProto(&created_torrent)
						new_torrents = append(new_torrents, &converted_torrent)
					} else {
						has_existing = true
					}
				}

				// Send the new torrents on the stream
				if callback != nil {
					if err := (*callback)(&pb.ScrapeResponse{
						MsDuration: uint32(time.Since(start).Milliseconds()),
						Torrents:   new_torrents,
					}); err != nil {
						return err
					}
				} else {
					log.Println("Scraped", len(new_torrents), "new torrents in", uint32(time.Since(start).Milliseconds()), "ms")
				}

				// Update NextPage to loop or complete the job
				page = page_result.NextPage
				if page == 0 || has_existing {
					break
				}
				time.Sleep(time.Second)
			}
			if err != nil {
				break
			}
		}

		// Handle timeout errors and skip to the next site
		if errors.Is(err, http.ErrHandlerTimeout) {
			err = nil
			continue
		} else {
			return err
		}
	}

	log.Println("Done ScrapeLatest")

	return nil
}

func (s *ScrapperServer) ScrapeLatest(request *pb.ScrapeLatestRequest, out pb.ScrapperService_ScrapeLatestServer) error {
	callback := func(response *pb.ScrapeResponse) error {
		return out.Send(response)
	}
	return DoScrapeLatest(&callback)
}
