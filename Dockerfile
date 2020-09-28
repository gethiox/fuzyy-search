# build environment
FROM golang:1.15.2-buster as builder

ADD src /app/src
WORKDIR /app/src

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/api/*.go

# run environment
FROM alpine:3.12.0

WORKDIR /root
COPY --from=builder /app/src/app .
CMD ["/root/app"]