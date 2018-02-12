FROM library/golang:1.9

RUN go get -u github.com/spf13/cobra
RUN go get -u github.com/spf13/viper
RUN go get -u github.com/go-sql-driver/mysql
# FOR win cross compilation
RUN go get -u github.com/inconshreveable/mousetrap
# To cross compile in general...
RUN go get -u github.com/mitchellh/gox

COPY . /go/src/github.com/odino/docsql

WORKDIR /go/src/github.com/odino/docsql
