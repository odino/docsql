package csv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io/ioutil"
)

// GetColumns | WARNING: highly inefficient software ahead!
// Here we read the whole TSV in order to get the
// list of columns -- this is not needed at all (it will
// consume more memory for the lolz) but, realistically,
// it's midnight and I got no jam left...
func GetColumns(filename string) ([]string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(bytes.NewReader(b))
	r.Comma = '\t'
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, errors.New("The spreadsheet seems to be empty. Are you sure it has rows?")
	}

	return records[0], nil
}
