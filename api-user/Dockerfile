FROM golang:1.18.3-alpine3.16

# https://levelup.gitconnected.com/docker-for-go-development-a27141f36ba9

# copy api files
WORKDIR /app
COPY ./ /app

# install api dependencies
RUN go mod download

# install CompileDaemon (hot reload)
RUN go get github.com/githubnemo/CompileDaemon
RUN go install -v github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main