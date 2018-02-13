FROM library/golang:1.9

RUN go get -u github.com/spf13/cobra \
            github.com/spf13/viper \
            github.com/go-sql-driver/mysql \
            # FOR win cross compilation
            github.com/inconshreveable/mousetrap \
            # To cross compile in general...
            github.com/mitchellh/gox

COPY . /go/src/github.com/odino/docsql

WORKDIR /go/src/github.com/odino/docsql
