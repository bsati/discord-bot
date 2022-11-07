FROM golang:1.19.3-alpine3.16

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build /app/cmd/migrate/migrate.go

CMD ["/app/migrate"]