package sites

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/golang/protobuf/ptypes"
	pb "github.com/trixky/hypertube/api-scrapper/proto"
	st "github.com/trixky/hypertube/api-scrapper/scrapper"
)

var URLS st.Urls = st.Urls{
	Movies: "http://www.legittorrents.info/index.php?page=torrents&active=1&category=1&order=3&by=2&pages=$page",
	Shows:  "http://www.legittorrents.info/index.php?page=torrents&active=1&category=13&order=3&by=2&pages=$page",
}

func scrapeList(page_type string, page uint32) (page_result st.ScrapperPageResult, err error) {
	torrents := make([]*pb.UnprocessedTorrent, 0, 20)
	c := colly.NewCollector(colly.IgnoreRobotsTxt())

	var category pb.MediaCategory
	if page_type == "movies" {
		category = pb.MediaCategory_CATEGORY_MOVIE
	} else {
		category = pb.MediaCategory_CATEGORY_SERIE
	}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("#bodyarea table table.lista table table.lista > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(index int, el *colly.HTMLElement) {
			if index > 0 {
				name := el.ChildText("td:nth-child(2) a")
				id_parts := strings.Split(el.ChildAttr("td:nth-child(2) a", "href"), "=")
				id := id_parts[len(id_parts)-1]
				seed64, _ := strconv.ParseInt(el.ChildText("td:nth-child(5)"), 10, 32)
				seed := int32(seed64)
				leech64, _ := strconv.ParseInt(el.ChildText("td:nth-child(6)"), 10, 32)
				leech := int32(leech64)
				// LegitTorrents has a single date format
				layout := "02/01/2006" // -- dd/mm/yyyy
				upload_time_date := el.ChildText("td:nth-child(4)")
				upload_time, _ := time.Parse(layout, upload_time_date)
				upload_timestamp, _ := ptypes.TimestampProto(upload_time)
				filename := strings.ReplaceAll(name, " ", "+")
				torrent_url := fmt.Sprintf("http://www.legittorrents.info/download.php?id=%s&f=%s.torrent", id, filename)

				torrents = append(torrents, &pb.UnprocessedTorrent{
					Name:       name,
					FullUrl:    fmt.Sprintf("http://www.legittorrents.info/index.php?page=torrent-details&id=%s", id),
					Type:       category,
					Seed:       &seed,
					Leech:      &leech,
					UploadTime: upload_timestamp,
					TorrentUrl: &torrent_url,
				})
			}
		})
	})

	c.OnHTML("form[name='change_pagepages']", func(h *colly.HTMLElement) {
		current_page := h.DOM.Find(".pagercurrent")
		if next_page := current_page.Next(); next_page != nil {
			page_value, _ := strconv.ParseUint(next_page.Text(), 10, 32)
			page_result.NextPage = uint32(page_value)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("error %v\n", err)
	})

	var url string
	if page_type == "movies" {
		url = URLS.Movies
	} else {
		url = URLS.Shows
	}

	if error := c.Visit(strings.Replace(url, "$page", fmt.Sprint(page), 1)); error != nil {
		fmt.Printf("error %v\n", error)
		return page_result, error
	}

	page_result.Torrents = torrents
	return
}

func scrapeSingle(id string) (torrent pb.UnprocessedTorrent, err error) {
	return
}

var Scrapper = st.Scrapper{
	ScrapeList:   scrapeList,
	ScrapeSingle: scrapeSingle,
}
