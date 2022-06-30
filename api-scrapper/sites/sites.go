package sites

import (
	st "github.com/trixky/hypertube/api-scrapper/scrapper"
	lt "github.com/trixky/hypertube/api-scrapper/sites/LegitTorrents"
)

var Scrappers = []st.Scrapper{
	lt.Scrapper,
}
