FROM golang:1.23.5-alpine AS builder

WORKDIR /usr/local/go/src/github.com/neo9/mongodb-backups

COPY . ./

# Descargar dependencias y generar el directorio vendor
RUN go mod tidy && go mod vendor

# Copiar el código fuente
COPY cmd ./cmd
COPY pkg ./pkg

ENV CGO_ENABLED=0

RUN cd cmd && go build -o /tmp/mongodb-backups

# Use an Alpine base image
FROM alpine:3.21.2

# Add the edge community repository
RUN echo 'http://dl-cdn.alpinelinux.org/alpine/edge/community' >> /etc/apk/repositories

# Update package list and install mongodb-tools
RUN apk update && \
    apk add --no-cache mongodb-tools

COPY --from=builder /tmp/mongodb-backups /bin/mongodb-backups

# Verify installation
RUN mongodump --version

CMD ["mongodb-backups"]
