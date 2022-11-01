FROM golang:1.19.2-alpine3.16 AS builder
WORKDIR /build
COPY . .
RUN go build -o daemon ./cmd/rde-daemon

FROM alpine:3.16
WORKDIR /daemon
COPY --from=builder /build/daemon daemon
RUN mkdir /var/rde /project
ENTRYPOINT [ "./daemon" ]