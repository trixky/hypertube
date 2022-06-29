package sites

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	pb "github.com/trixky/hypertube/api-scrapper/proto"
	st "github.com/trixky/hypertube/api-scrapper/scrapper"
)

var URLS st.Urls = st.Urls{
	Movies: "http://www.legittorrents.info/index.php?page=torrents&active=1&category=1&order=3&by=2&pages=$page",
	Shows:  "http://www.legittorrents.info/index.php?page=torrents&active=1&category=13&order=3&by=2&pages=$page",
}

func scrapeList(page_type string, page uint32) (page_result st.ScrapperPageResult, err error) {
	torrents := make([]*pb.Torrent, 0, 20)
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
				seed64, _ := strconv.ParseUint(el.ChildText("td:nth-child(5)"), 10, 32)
				seed := uint32(seed64)
				leech64, _ := strconv.ParseUint(el.ChildText("td:nth-child(6)"), 10, 32)
				leech := uint32(leech64)
				size := ""
				// TODO Convert to DateTime
				// upload_time := el.ChildText("td:nth-child(4)")
				// TODO Name to filename (remove spaces and disallowed characters)
				torrent_url := fmt.Sprintf("http://www.legittorrents.info/download.php?id=%s&f=%s.torrent", id, name)

				torrents = append(torrents, &pb.Torrent{
					Id:         0,
					Name:       name,
					FullUrl:    fmt.Sprintf("http://www.legittorrents.info/index.php?page=torrent-details&id=%s", id),
					Type:       category,
					Seed:       &seed,
					Leech:      &leech,
					Size:       &size,
					UploadTime: nil,
					TorrentUrl: torrent_url,
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

func scrapeSingle(id string) (torrent pb.Torrent, err error) {
	return
}

var Scrapper = st.Scrapper{
	ScrapeList:   scrapeList,
	ScrapeSingle: scrapeSingle,
}
