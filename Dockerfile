FROM golang:latest as builder
WORKDIR /go/src/github.com/jdnielss/gonar
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gonar .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/jdnielss/gonar/gonar .
ENTRYPOINT ["./gonar"]
