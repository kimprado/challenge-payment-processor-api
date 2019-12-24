# compile stage
FROM golang:1.13.4 AS build

ARG workbuild=/usr/dist

COPY    go.mod /src/
COPY    go.sum /src/
WORKDIR /src/
RUN     go mod download
RUN     go get github.com/google/wire/cmd/wire@v0.4.0

EXPOSE  3000
