# Use Multistate build

ARG APP_PORT=8080

FROM golang:1.18-alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download


FROM golang:1.18-alpine as builder
RUN apk add --update --no-cache ca-certificates
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd/app

FROM scratch
EXPOSE $APP_PORT
COPY --from=builder /app/config /config
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
ENTRYPOINT ["/app"]
