# create build image
FROM golang:alpine AS builder

# set golang env
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# cd to work dir
WORKDIR /work

# copy source code file to WORKDIR
COPY . ./

# build
RUN go build -o entry .

# create running image
FROM debian:bullseye-slim As latest
ARG RUN_ENV

# copy config and sql file
COPY config.${RUN_ENV}.json ./
COPY init.sql ./

# copy static files
COPY static ./static

# copy execuable mgtbe file from builder image to current dir
COPY --from=builder /work/entry /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

