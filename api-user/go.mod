module github.com/trixky/hypertube/api-user

go 1.18

require github.com/trixky/hypertube/.shared v0.0.0

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.3
	github.com/lib/pq v1.10.6
	google.golang.org/genproto v0.0.0-20220630174209-ad1d48641aa7
	google.golang.org/grpc v1.47.0
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.19.0 // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/trixky/hypertube/.shared v0.0.0 => ../.shared
