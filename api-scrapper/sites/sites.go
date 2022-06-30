package sites

import (
	st "github.com/trixky/hypertube/api-scrapper/scrapper"
	lx "github.com/trixky/hypertube/api-scrapper/sites/L337x"
	lt "github.com/trixky/hypertube/api-scrapper/sites/LegitTorrents"
)

var Scrappers = []st.Scrapper{
	lt.Scrapper,
	lx.Scrapper,
}
