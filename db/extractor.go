package db

import (
	"regexp"
	"strconv"
)

type connection struct {
	user     string
	pass     string
	host     string
	database string
	port     int64
}

func extract(conn string) (*connection, error) {

	r, err := regexp.Compile("^([a-zA-Z0-9]+)[:]([a-zA-Z0-9])?[^(]+[(]([^:]+)[:]([0-9]+)[)][/]([^?]+).*$")
	if err != nil {
		return nil, err
	}
	match := r.FindStringSubmatch(conn)
	port, err := strconv.ParseInt(match[4], 0, 64)
	if err != nil {
		return nil, err
	}
	return &connection{
		user:     match[1],
		pass:     match[2],
		host:     match[3],
		port:     port,
		database: match[5],
	}, nil
}
