FROM golang:1.12.7-alpine3.10 as builder

WORKDIR /usr/local/go/src/github.com/neo9/mongodb-backups

RUN apk add git &&  export GO111MODULE=on

ADD go.mod ./
ADD go.sum ./

RUN go mod vendor

ADD cmd ./cmd
ADD pkg ./pkg
RUN ls -lah

RUN cd cmd  && go build -o /tmp/mongodb-backups

FROM alpine:3.10

COPY --from=builder /tmp/mongodb-backups /bin/mongodb-backup

CMD mongodb-backup
