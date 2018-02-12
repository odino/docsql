package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/odino/docsql/csv"
	"github.com/odino/docsql/db"
	"github.com/odino/docsql/gdocs"
	"github.com/odino/docsql/util"

	"github.com/spf13/cobra"
	config "github.com/spf13/viper"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "docsql",
	Short: "docsql imports a spreadsheet hosted on Google Docs to a MySQL table",
	Long: `A tool to import Google Docs' spreadhsheets into a MySQL table.
Have a look at https://github.com/odino/docsql for more info`,
	Run: func(cmd *cobra.Command, args []string) {
		tablename := config.GetString("table") + "_" + strconv.FormatInt(time.Now().UnixNano(), 10)
		filename := tablename + ".csv" // What about ioutil.TempFile?
		defer os.Remove(filename)
		// Download doc
		err := gdocs.Download(config.GetString("doc"), filename, config.GetInt64("timeout"))
		util.Check(err)

		columns, err := csv.GetColumns(filename)
		util.Check(err)

		err = db.CreateTable(config.GetString("connection"), tablename, columns)
		util.Check(err)

		err = db.LoadData(config.GetString("connection"), tablename, filename)
		util.Check(err)

		err = db.RenameTables(config.GetString("connection"), tablename, config.GetString("table"))
		util.Check(err)

		err = db.DeleteArchiveTables(config.GetString("connection"), config.GetString("table"), config.GetInt("keep"))
		if err != nil {
			log.Printf("Unable to delete archived tables")
		}

		log.Printf("All done")
		// cleanup
		// minrows
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	config.AutomaticEnv()
	rootCmd.Flags().StringP("doc", "d", "", "Url of the spreadsheet to sync to MySQL")
	config.BindPFlag("doc", rootCmd.Flags().Lookup("doc"))
	rootCmd.MarkFlagRequired("doc")

	rootCmd.Flags().StringP("table", "t", "", "Name of the table you want to dump the spreadsheet into")
	config.BindPFlag("table", rootCmd.Flags().Lookup("table"))
	rootCmd.MarkFlagRequired("table")

	rootCmd.Flags().StringP("connection", "c", "", "Connection string for MySQL (eg. 'root:root@tcp(localhost:3306)/test_database?charset=utf8')")
	config.BindPFlag("connection", rootCmd.Flags().Lookup("connection"))

	// This makes it easy to pass the MySQL credentials via environment
	if config.GetString("connection") == "" {
		rootCmd.MarkFlagRequired("connection")
	}

	rootCmd.Flags().Int64P("timeout", "T", 5, "Timeout, in seconds, to download the Google Doc")
	config.BindPFlag("timeout", rootCmd.Flags().Lookup("timeout"))

	rootCmd.Flags().IntP("keep", "k", 10, "Tables to keep in MySQL, generally for historical purpose or to be able to rollback")
	config.BindPFlag("keep", rootCmd.Flags().Lookup("keep"))
}
