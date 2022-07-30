FROM golang:1.18.3-alpine3.16 as builder

WORKDIR /hypertube/app
COPY ./api-user /hypertube/app
COPY ./.shared /hypertube/.shared

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o /api

FROM golang:1.18.3-alpine3.16

RUN addgroup -S hypertube && adduser -S hypertube -G hypertube

RUN mkdir -p /hypertube/app/storage
RUN chown -R hypertube:hypertube /hypertube/app/storage
RUN mkdir -p /api
RUN chown -R hypertube:hypertube /api

RUN mkdir /data
COPY --from=builder /api .

USER hypertube
ENTRYPOINT ./api
