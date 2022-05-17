# Dockerfile to build the Change Event (or CloudEvent) interceptor
FROM golang:1.18-bullseye as builder

WORKDIR /go/ce-interceptor
COPY . .
RUN go build .

FROM gcr.io/distroless/base:nonroot
COPY --from=builder /go/ce-interceptor/ce-interceptor /ce-interceptor

EXPOSE 8080

ENTRYPOINT ["/ce-interceptor"]
