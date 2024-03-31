FROM golang:1.22 as builder

ENV CGO_ENABLED=0

WORKDIR /app
RUN apt-get update && apt-get install -y gcc musl-dev sqlite3 libsqlite3-dev

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
COPY . .
RUN cd ./app && go build -o /app/esport-syncer .

FROM alpine:3.19

COPY --from=builder /app/esport-syncer /srv/esport-syncer
WORKDIR /srv

ENTRYPOINT ["/srv/esport-syncer"]

#ENTRYPOINT ["tail"]
#CMD ["-f","/dev/null"]
