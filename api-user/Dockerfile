FROM golang:1.18.3-alpine3.16 as builder

RUN addgroup -S hypertube && adduser -S hypertube -G hypertube

WORKDIR /hypertube/app
COPY ./api-user /hypertube/app
COPY ./.shared /hypertube/.shared

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o /api

FROM golang:1.18.3-alpine3.16

RUN mkdir /data
COPY --from=builder /api .

USER hypertube
ENTRYPOINT ./api
