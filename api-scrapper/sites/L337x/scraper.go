package sites

import (
	"fmt"
	"log"
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
	Movies: "https://www.1377x.to/sort-cat/Movies/time/desc/$page/",
	Shows:  "https://www.1377x.to/sort-cat/TV/time/desc/$page/",
}

var time_prefixes []string = []string{
	"th", "st", "nd", "rd",
}

// 1337x has 3 date formats
// test: https://go.dev/play/p/koVP_R5IaEi
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
	c := colly.NewCollector(colly.IgnoreRobotsTxt(), colly.UserAgent("Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:101.0) Gecko/20100101 Firefox/101.0"))
	re := regexp.MustCompile("\\d+$")

	var category pb.MediaCategory
	if page_type == "movies" {
		category = pb.MediaCategory_CATEGORY_MOVIE
	} else {
		category = pb.MediaCategory_CATEGORY_SERIE
	}

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnHTML(".featured-list table.table-list.table tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			// Try to use the HTML content of the name node since Colly add ellipsis somehow
			name := ""
			name_node := el.DOM.Find("td.name a[href^='/torrent/']")
			if name_node != nil {
				name, _ = name_node.Html()
			} else {
				name = el.ChildText("td.name a[href^='/torrent/']")
			}

			relative_href := el.ChildAttr("td.name a[href^='/torrent/']", "href")
			if relative_href == "" || relative_href == "/" {
				log.Println("found empty relative url in", e.Request.URL, "for", name)
				return
			}
			full_url := "https://1337x.to" + relative_href

			seed64, _ := strconv.ParseInt(el.ChildText("td.seeds"), 10, 32)
			seed := int32(seed64)

			leech64, _ := strconv.ParseInt(el.ChildText("td.leeches"), 10, 32)
			leech := int32(leech64)

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
			})
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("error", err)
	})

	var url string
	if page_type == "movies" {
		url = URLS.Movies
	} else {
		url = URLS.Shows
	}

	if err := c.Visit(strings.Replace(url, "$page", fmt.Sprint(page), 1)); err != nil {
		log.Println("error", err)
		return page_result, err
	}

	c.Wait()

	if page < 150 {
		page_result.NextPage = page + 1
	}
	page_result.Torrents = torrents
	return
}

func scrapeSingle(torrent *pb.UnprocessedTorrent) error {
	c := colly.NewCollector(colly.IgnoreRobotsTxt())

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnHTML(".torrent-detail-page", func(e *colly.HTMLElement) {
		description := e.DOM.Find("#description")
		if description != nil {
			description, _ := description.Html()
			torrent.DescriptionHtml = &description
			imdb_id_match := st.IMDBre.FindStringSubmatch(description)
			if len(imdb_id_match) == 2 {
				torrent.ImdbId = &imdb_id_match[1]
			}
		}

		magnet := e.ChildAttr("a[href^='magnet:']", "href")
		if magnet != "" {
			torrent.Magnet = &magnet
		}

		torrent_url := e.ChildAttr(".dropdown-menu > li > a", "href")
		torrent.TorrentUrl = &torrent_url

		// test: https://go.dev/play/p/4QLPAxQkcuM
		re := regexp.MustCompile("(.+)\\s+\\((.+)\\)\\s*$")
		e.ForEach("#files > ul > li", func(index int, h *colly.HTMLElement) {
			name := h.Text
			matches := re.FindStringSubmatch(name)
			if len(matches) == 3 {
				path, name := st.ExtractPath(name)
				torrent.Files = append(torrent.Files, &pb.TorrentFile{
					Name: name,
					Path: path,
					Size: &matches[2],
				})
			}
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("error", err)
	})

	if err := c.Visit(torrent.FullUrl); err != nil {
		log.Println("error", err)
		return err
	}

	c.Wait()

	return nil
}

var Scrapper = st.Scrapper{
	ScrapeList:   scrapeList,
	ScrapeSingle: scrapeSingle,
}
