FROM golang:1.23-alpine as build

WORKDIR /app

COPY go.mod /go.sum /app/
RUN go mod download

COPY . /app/

RUN CGO_ENABLED=0 go build -o /webhook ./cmd

FROM alpine:3.20

COPY --from=build /webhook /usr/local/bin/webhook
RUN chmod +x /usr/local/bin/webhook

ENTRYPOINT ["webhook"]
