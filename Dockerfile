FROM golang:1.14
RUN apk add bash mysql-client

RUN mkdir -p /home/go/api && chown -R go:go /home/go/api

WORKDIR /home/go/api

COPY . .

USER go

RUN GOOS=linux go build -ldflags="-s -w"

COPY --chown=go:go . .

EXPOSE 3333

ENTRYPOINT ["./main"]