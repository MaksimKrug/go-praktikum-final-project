FROM golang:1.24.13-bookworm AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY pkg ./pkg
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/scheduler .

FROM debian:bookworm-slim

WORKDIR /app
RUN mkdir -p /data && chown 10001:10001 /data
COPY --from=build /out/scheduler ./scheduler
COPY --from=build /usr/local/go/lib/time/zoneinfo.zip ./zoneinfo.zip
COPY web ./web

ENV TODO_PORT=7540 \
    TODO_DBFILE=/data/scheduler.db \
    TZ=UTC \
    ZONEINFO=/app/zoneinfo.zip

USER 10001:10001
EXPOSE 7540
ENTRYPOINT ["/app/scheduler"]
