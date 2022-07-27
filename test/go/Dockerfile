FROM postgres:14.4-alpine

# ------------------------------- copy test
COPY ./test/go /hypertube/test/go

# ------------------------------- redis install
RUN apk add redis

# ------------------------------- postgres configuration
COPY ./postgres /hypertube/postgres
COPY ./postgres/init.sql /docker-entrypoint-initdb.d/

# ------------------------------- install go
COPY --from=golang:1.18.3-alpine3.16 /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"
RUN apk add gcc musl-dev 

# ------------------------------- copy project
WORKDIR /hypertube/
COPY ./.shared /hypertube/.shared
COPY ./api-auth /hypertube/api-auth
COPY ./api-user /hypertube/api-user

# it may be interesting to rebuild the image without
# cache to take advantage of the pre-download of go dependencies
# ------------------------------- pre-download dependencies
# .shared
WORKDIR /hypertube/.shared
RUN go mod download
# api-auth
WORKDIR /hypertube/api-auth
RUN go mod download
# api-user
WORKDIR /hypertube/api-user
RUN go mod download

WORKDIR /hypertube

ENTRYPOINT /bin/bash /hypertube/test/go/test.sh