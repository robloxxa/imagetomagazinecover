FROM golang:1.21.0-alpine3.18 as build
WORKDIR /app
COPY go.mod go.sum ./
COPY *.go .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /kachenokmagazinebot

FROM chromedp/headless-shell:latest

COPY ca-bundle.crt /etc/ssl/certs/ca-bundle.crt
COPY ca-bundle.trust.crt /etc/ssl/certs/ca-bundle.trust.crt
COPY --from=build /kachenokmagazinebot /kachenokmagazinebot
COPY static/ /static/

RUN apt update
RUN apt install dumb-init

ENTRYPOINT ["dumb-init", "--"]

CMD [ "/kachenokmagazinebot", "-PORT $PORT", "-MAX_WORKERS ${MAX_WORKERS}" ]

