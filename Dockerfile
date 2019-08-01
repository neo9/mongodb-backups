FROM golang:1.12.7-alpine3.10 as builder

WORKDIR /usr/local/go/src/github.com/neo9/mongodb-backups

RUN apk add git &&  export GO111MODULE=on

ADD go.mod ./
ADD go.sum ./

RUN go mod vendor

ADD cmd ./cmd
ADD pkg ./pkg

RUN cd cmd  && go build -o /tmp/mongodb-backups

FROM alpine:3.9

COPY --from=builder /tmp/mongodb-backups /bin/mongodb-backup
ENV MONGODB_TOOLS_VERSION 4.0.5-r0

RUN apk add --no-cache mongodb-tools=${MONGODB_TOOLS_VERSION}

CMD mongodb-backup
