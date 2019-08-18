FROM library/golang:1.12-alpine

RUN apk add --update git bash
RUN go get -u github.com/spf13/cobra \
            github.com/spf13/viper \
            github.com/go-sql-driver/mysql \
            github.com/mitchellh/gox

COPY . /src

WORKDIR /src
RUN go mod vendor
