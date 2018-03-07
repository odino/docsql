package db

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Dummy function used to sanitize column names
// and avoid funny characters like ` to break
// the CREATE TABLE statement.
func sanitize(s string, fallback string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9_]+")
	if err != nil {
		return "___ERROR___"
	}

	s = reg.ReplaceAllString(s, "")

	if s == "" {
		return fallback
	}

	return s
}

func getCreateTableQuery(tablename string, columns []string) string {
	cols := []string{}

	for i, col := range columns {
		cols = append(cols, fmt.Sprintf("`%s` VARCHAR(255) NOT NULL COLLATE 'utf8_unicode_ci',", sanitize(col, strconv.Itoa(i))))
	}

	return fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
      %s
    	docsql_id INT(11) NOT NULL AUTO_INCREMENT,
    	docsql_created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    	PRIMARY KEY (docsql_id)
    )
    COLLATE='utf8_unicode_ci'
    ENGINE=InnoDB;
    `, tablename, strings.Join(cols, " "))
}

func CreateTable(conn string, tablename string, columns []string) error {
	log.Printf("Connecting to MySQL...")
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return err
	}

	defer db.Close()

	log.Printf("Creating table '%s'...", tablename)
	_, err = db.Query(getCreateTableQuery(tablename, columns))
	if err != nil {
		return err
	}

	return nil
}

func LoadData(conn string, tablename string, filename string) error {
	log.Printf("Connecting to MySQL...")
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return err
	}

	defer db.Close()

	log.Printf("Loading data into '%s'...", tablename)
	_, err = db.Query(fmt.Sprintf("LOAD DATA LOCAL INFILE './%s' INTO TABLE %s FIELDS TERMINATED BY '\t' LINES TERMINATED BY '\r\n' IGNORE 1 LINES", filename, tablename))
	if err != nil {
		return err
	}

	return nil
}

// We rename / swap tables together so that the operation is atomic.
// See: https://stackoverflow.com/a/34391961/934439
func RenameTables(conn string, newtable string, oldtable string) error {
	log.Printf("Connecting to MySQL...")
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return err
	}

	defer db.Close()

	err = CreateTable(conn, oldtable, []string{})
	if err != nil {
		return err
	}

	log.Printf("Swapping '%s' with '%s'", oldtable, newtable)
	_, err = db.Query(fmt.Sprintf("RENAME TABLE %s TO %s_archive, %s TO %s", oldtable, newtable, newtable, oldtable))
	if err != nil {
		return err
	}

	return nil
}

// We need to make sure that everytime docsql runs it takes care of removing
// old swapped tables, else they'll keep accumulating. Here we fail-safe, in the
// sense that if there's any issue DROPping the old tables we simply log it, but
// don't throw any error.
func DeleteArchiveTables(conn string, table string, keep int) error {
	log.Printf("Connecting to MySQL...")
	extract, err := extract(conn)
	if err != nil {
		return err
	}
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Printf("Clearing old tables...")
	tables, err := db.Query(fmt.Sprintf("SELECT table_name FROM information_schema.tables where table_schema='%s' AND table_name LIKE \"%s_%%_archive\" ORDER BY table_name DESC LIMIT 10000 OFFSET %d", extract.database, table, keep))
	if err != nil {
		return err
	}

	defer tables.Close()
	for tables.Next() {
		var name string
		err = tables.Scan(&name)
		_, err := db.Query(fmt.Sprintf("DROP TABLE %s", name))
		if err != nil {
			log.Printf("Unable to DROP table '%s'", name)
		}
	}

	return tables.Err()
}
