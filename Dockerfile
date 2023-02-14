FROM golang:1.20-alpine3.17 AS build

WORKDIR /src
COPY go.* .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/marymagazinebot .

FROM alpine:3.17.2
COPY -- from=build /marymagazinebot /marymagazinebot