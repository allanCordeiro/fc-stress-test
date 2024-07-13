FROM golang:1.22-alpine AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /fc-stress-test main.go

FROM scratch
COPY --from=builder /fc-stress-test /fc-stress-test
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


ENTRYPOINT ["/fc-stress-test"]