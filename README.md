# docsql

A tool to import spreadsheets hosted on Google Docs to a MySQL table.

## Usage

Grab a binary from the [releases' page](https://github.com/odino/docsql/releases)
and start having some fun:

``` bash
$ docsql \
--doc "https://docs.google.com/spreadsheets/d/1vyVxaYgfZ2Tka7reg4whg99kRlWqpg6cKvEa1QFArZI/export?format=tsv" \
--table my_sample \
--connection "root:@tcp(localhost:3308)/test?charset=utf8&allowAllFiles=true"

2018/02/12 23:20:31 Downloading https://docs.google.com/spreadsheets/d/1vyVxaYgfZ2Tka7reg4whg99kRlWqpg6cKvEa1QFArZI/export?format=tsv ...
2018/02/12 23:20:32 Doc downloaded in my_sample_1518463231621589126.csv
2018/02/12 23:20:32 Connecting to MySQL...
2018/02/12 23:20:32 Creating table 'my_sample_1518463231621589126'...
2018/02/12 23:20:32 Connecting to MySQL...
2018/02/12 23:20:32 Loading data into 'my_sample_1518463231621589126'...
2018/02/12 23:20:32 Connecting to MySQL...
2018/02/12 23:20:32 Swapping 'my_sample' with 'my_sample_1518463231621589126'
2018/02/12 23:20:32 Connecting to MySQL...
2018/02/12 23:20:32 Creating table 'my_sample'...
2018/02/12 23:20:32 Connecting to MySQL...
2018/02/12 23:20:32 Clearing old tables...
2018/02/12 23:20:32 All done
```

![example](https://raw.githubusercontent.com/odino/docsql/master/images/docsql.png?raw=true)

## Advanced

### Spreadsheet

Your spreadsheet will need to be shared publicly (*anyone with the link can access*),
and the URL you need to feed to docsql takes the form of `https://docs.google.com/spreadsheets/d/$DOCID/export?format=tsv`
where `$DOCID` is the unique ID of the Google Doc.

By default, docsql will download the first sheet in the doc, but if you need to
import other sheets you can simply append the `gid` of the sheet at the end of the URL
(`https://docs.google.com/spreadsheets/d/$DOCID/export?format=tsv&gid=$GID`).

Please note that the export format must be **tsv** because, well, [it's just
easier than csv](https://en.wikipedia.org/wiki/Tab-separated_values)

### MySQL

Instead of passing the connection string to MySQL as a flag you can export it as
environment variable -- this makes sure you don't leave credentials on the CLI:

``` bash
$ export $CONNECTION=...

$ docsql \
--doc "https://docs.google.com/spreadsheets/d/1vyVxaYgfZ2Tka7reg4whg99kRlWqpg6cKvEa1QFArZI/export?format=tsv" \
--table my_sample  

2018/02/12 23:27:30 Downloading https://docs.google.com/spreadsheets/d/1vyVxaYgfZ2Tka7reg4whg99kRlWqpg6cKvEa1QFArZI/export?format=tsv ...
2018/02/12 23:27:33 Doc downloaded in my_sample_1518463650997899367.csv
2018/02/12 23:27:33 Connecting to MySQL...
2018/02/12 23:27:33 Creating table 'my_sample_1518463650997899367'...
2018/02/12 23:27:33 Connecting to MySQL...
2018/02/12 23:27:33 Loading data into 'my_sample_1518463650997899367'...
2018/02/12 23:27:33 Connecting to MySQL...
2018/02/12 23:27:33 Swapping 'my_sample' with 'my_sample_1518463650997899367'
2018/02/12 23:27:33 Connecting to MySQL...
2018/02/12 23:27:33 Creating table 'my_sample'...
2018/02/12 23:27:33 Connecting to MySQL...
2018/02/12 23:27:33 Clearing old tables...
2018/02/12 23:27:33 All done
```

Be aware that `LOAD DATA LOCAL INFILE` must be available on the MySQL server,
and you will need to end your connection string with `allowAllFiles=true` so that
the Go MySQL driver is allowed to process local files.

### Keeping old tables

docsql is (probably) meant to run as a cron, or everytime you make an update to
your spreadsheet -- whenever it runs, it nukes the previous version of the output
table and imports the new contents of the spreadsheet.

You can customize how many (old) tables to keep with the `--keep` flag. For example,
`docsql ... --keep 5` will keep 5 version of the old table in MySQL:

``` bash
mysql> SHOW TABLES;
+---------------------------------------+
| Tables_in_test                        |
+---------------------------------------+
| my_sample                             |
| my_sample_1518463163413558194_archive |
| my_sample_1518463168405819860_archive |
| my_sample_1518463173716215291_archive |
| my_sample_1518463231621589126_archive |
| my_sample_1518463650997899367_archive |
+---------------------------------------+
```

### Table structure

docsql will make a few opinionated assumptions for you:

* all fields in the table are `VARCHAR(255)`
* it creates an `docsql_id` field used as a primary key
* it adds an `docsql_created_at` with the timestamp when the rows were loaded into the table
* will sanitize column names (taken from the spreadsheet) filtering out non alphanumeric characters

There are plans to make all of these configurable in the future through flags...
...PRs are more than welcome!

### Other stuff?

It might be a good idea to run `docsql --help` to have a look at what's available.

## Contributing

docsql is being developed through docker because... ...well, don't always have
the Go toolchain with me!

Anyhow, it should be fairly straighforward to get running:

* `make build_docker`, will build the docker container used to develop

Feel free to rant or, even better, fix some of my crappy code through a [pull request](https://github.com/odino/docsql/pulls)!
