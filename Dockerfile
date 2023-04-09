FROM golang:1.20 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o esport-syncer github.com/pkarpovich/esport-syncer/app

FROM alpine:3.17

WORKDIR /
COPY --from=builder /app/esport-syncer /esport-syncer

ENTRYPOINT ["/esport-syncer"]

#ENTRYPOINT ["tail"]
#CMD ["-f","/dev/null"]
