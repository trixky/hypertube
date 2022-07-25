module github.com/trixky/hypertube/api-auth

go 1.18

require github.com/trixky/hypertube/.shared v0.0.0

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.0
	github.com/lib/pq v1.10.6
	golang.org/x/oauth2 v0.0.0-20220722155238-128564f6959c
	google.golang.org/genproto v0.0.0-20220719170305-83ca9fad585f
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.0
)

require (
	cloud.google.com/go/compute v1.7.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/onsi/ginkgo v1.12.1 // indirect
	github.com/onsi/gomega v1.10.3 // indirect
	golang.org/x/net v0.0.0-20220624214902-1bab6f366d9e // indirect
	golang.org/x/sys v0.0.0-20220610221304-9f5ed59c137d // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/trixky/hypertube/.shared v0.0.0 => ../.shared
