# Dockerfile to build out the image for the service to retrieve data from Redis

FROM golang:1.18-bullseye as builder

WORKDIR /go/redis-retriever
COPY . .
RUN go build .

FROM gcr.io/distroless/base:nonroot
COPY --from=builder /go/redis-retriever/redis-retriever /redis-retriever

EXPOSE 8080
ENTRYPOINT ["/redis-retriever"]
