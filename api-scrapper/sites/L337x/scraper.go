package sites

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/golang/protobuf/ptypes"
	pb "github.com/trixky/hypertube/api-scrapper/proto"
	st "github.com/trixky/hypertube/api-scrapper/scrapper"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var URLS st.Urls = st.Urls{
	Movies: "https://1337x.to/cat/Movies/$page/",
	Shows:  "https://1337x.to/cat/TV/$page/",
}

var time_prefixes []string = []string{
	"th", "st", "nd", "rd",
}

// 1337x has 3 date formats
func parseDateTime(value string) (timestamp *timestamppb.Timestamp, err error) {
	layout := "3:04pm" // -- XX:XX{am,pm}
	upload_time, err := time.Parse(layout, value)
	if err == nil {
		today := time.Now()
		upload_time = upload_time.AddDate(today.Year(), int(today.Month()), today.Day())
	}
	if err != nil {
		for _, prefix := range time_prefixes {
			layout := "3pm Jan. 2" + prefix // -- X{am,pm} ShortMonth. day{st,nd,rd,th}
			upload_time, err = time.Parse(layout, value)
			if err == nil {
				today := time.Now()
				upload_time = upload_time.AddDate(today.Year(), 0, 0)
				break
			}
		}
	}
	if err != nil {
		for _, prefix := range time_prefixes {
			layout := "Jan. 2" + prefix + " '06" // -- ShortMonth. day{st,nd,rd,th} 'YearPrefix
			upload_time, err = time.Parse(layout, value)
			if err == nil {
				break
			}
		}
	}
	if err == nil {
		return ptypes.TimestampProto(upload_time)
	}
	return
}

func scrapeList(page_type string, page uint32) (page_result st.ScrapperPageResult, err error) {
	torrents := make([]*pb.UnprocessedTorrent, 0, 25)
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

	c.OnHTML(".featured-list table.table-list.table tbody ", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			name := el.ChildText("td.name a:last-child")
			relative_href := el.ChildAttr("td.name a:last-child", "href")
			full_url := "https://1337x.to" + relative_href
			seed64, _ := strconv.ParseUint(el.ChildText("td.seeds"), 10, 32)
			seed := uint32(seed64)
			leech64, _ := strconv.ParseUint(el.ChildText("td.leeches"), 10, 32)
			leech := uint32(leech64)
			re := regexp.MustCompile("\\d+$")
			size := el.ChildText("td.size")
			size = re.ReplaceAllString(size, "")
			upload_date_value := el.ChildText("td.coll-date")
			upload_timestamp, _ := parseDateTime(upload_date_value)

			torrents = append(torrents, &pb.UnprocessedTorrent{
				Name:       name,
				FullUrl:    full_url,
				Type:       category,
				Seed:       &seed,
				Leech:      &leech,
				Size:       &size,
				UploadTime: upload_timestamp,
				TorrentUrl: nil,
			})
		})
	})

	c.OnHTML("div.pagination", func(h *colly.HTMLElement) {
		if page < 150 {
			current_page := h.DOM.Find(".active")
			if next_page := current_page.Next(); next_page != nil {
				page_result.NextPage = page + 1
			}
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
