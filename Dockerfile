FROM golang:1.16.8-alpine3.14 AS builder

WORKDIR /usr/local/go/src/github.com/neo9/mongodb-backups

RUN apk add git &&  export GO111MODULE=on

ADD go.mod ./
ADD go.sum ./

RUN go mod vendor

ADD cmd ./cmd
ADD pkg ./pkg

ENV CGO_ENABLED=0

RUN cd cmd  && go build -o /tmp/mongodb-backups

FROM debian:10-slim

RUN apt-get update && apt-get install curl -y

COPY --from=builder /tmp/mongodb-backups /bin/mongodb-backups

RUN curl -o mongodb-tools.deb https://fastdl.mongodb.org/tools/db/mongodb-database-tools-debian10-x86_64-100.5.0.deb && apt install ./mongodb-tools.deb && rm mongodb-tools.deb
RUN rm -rf /var/cache/apt

CMD mongodb-backups
