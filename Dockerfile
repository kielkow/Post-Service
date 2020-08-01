FROM golang:1.14

RUN apk add --no-cache openssl bash mysql-client

RUN apk add --no-cache shadow

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

RUN usermod -u 1000 go

RUN mkdir -p /home/go/api && chown -R go:go /home/go/api

WORKDIR /home/go/api

COPY . .

USER go

RUN GOOS=linux go build -ldflags="-s -w"

COPY --chown=go:go . .

EXPOSE 3333

ENTRYPOINT ["./main"]