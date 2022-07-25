module github.com/trixky/hypertube/api-scrapper

go 1.18

require github.com/trixky/hypertube/.shared v0.0.0

require (
	github.com/go-co-op/gocron v1.15.1
	github.com/gocolly/colly v1.2.0
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.0
	github.com/lib/pq v1.10.6
	google.golang.org/genproto v0.0.0-20220722212130-b98a9ff5e252
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/PuerkitoBio/goquery v1.8.0 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/antchfx/htmlquery v1.2.5 // indirect
	github.com/antchfx/xmlquery v1.3.11 // indirect
	github.com/antchfx/xpath v1.2.1 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca // indirect
	github.com/temoto/robotstxt v1.1.2 // indirect
	golang.org/x/net v0.0.0-20220624214902-1bab6f366d9e // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220610221304-9f5ed59c137d // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.6.7 // indirect
)

replace github.com/trixky/hypertube/.shared v0.0.0 => ../.shared
